package frame

import (
	"context"
	"fmt"
	"github.com/Remember9/frame/config"
	"github.com/Remember9/frame/flag"
	"github.com/Remember9/frame/server"
	"github.com/Remember9/frame/signals"
	"github.com/Remember9/frame/util/xcycle"
	"github.com/Remember9/frame/util/xdefer"
	"github.com/Remember9/frame/util/xgo"
	"github.com/Remember9/frame/util/xstring"
	"github.com/Remember9/frame/worker"
	job "github.com/Remember9/frame/worker/xjob"
	"github.com/Remember9/frame/xlog"
	"go.uber.org/automaxprocs/maxprocs"
	"golang.org/x/sync/errgroup"
	"runtime"
	"sync"
)

type Application struct {
	name        string
	version     string
	env         string
	cycle       *xcycle.Cycle
	stopOnce    sync.Once
	initOnce    sync.Once
	startupOnce sync.Once
	afterStop   *xdefer.DeferStack
	beforeStop  *xdefer.DeferStack
	defers      []func() error
	servers     []server.Server
	workers     []worker.Worker
	jobs        map[string]job.Runner
}

func (app *Application) initialize() {
	app.initOnce.Do(func() {
		app.cycle = xcycle.NewCycle()
		app.afterStop = xdefer.NewStack()
		app.beforeStop = xdefer.NewStack()
		app.servers = make([]server.Server, 0)
		app.workers = make([]worker.Worker, 0)
		app.jobs = make(map[string]job.Runner)
	})
}

func (app *Application) startup() (err error) {
	app.startupOnce.Do(func() {
		err = xgo.SerialUntilError(
			app.printBanner,
			app.initConfig,
			app.initLogger,
			app.initFrame,
			app.printInitEnd,
		)()
	})
	return
}

func (app *Application) Startup(fns ...func() error) error {
	app.initialize()
	if err := app.startup(); err != nil {
		return err
	}
	return xgo.SerialUntilError(fns...)()
}

func (app *Application) Run() error {
	app.waitSignals()
	defer app.clean()

	if len(app.jobs) > 0 {
		_ = app.startJobs()
	} else {
		app.cycle.Run(app.startServers)
		app.cycle.Run(app.startWorkers)
		if err := <-app.cycle.Wait(); err != nil {
			xlog.Error("shutdown with error", xlog.FieldErr(err))
		}
	}

	xlog.Info("shutdown, bye!")
	return nil
}

func (app *Application) Stop() (err error) {
	app.stopOnce.Do(func() {
		app.beforeStop.Clean()

		for _, s := range app.servers {
			func(s server.Server) {
				app.cycle.Run(s.Stop)
			}(s)
		}
		for _, w := range app.workers {
			func(w worker.Worker) {
				app.cycle.Run(w.Stop)
			}(w)
		}

		select {
		case <-app.cycle.Done():
			app.cycle.Close()
		}
	})
	return
}

func (app *Application) BeforeStop(fns ...func() error) {
	app.initialize()
	if app.beforeStop == nil {
		app.beforeStop = xdefer.NewStack()
	}
	app.beforeStop.Push(fns...)
}

func (app *Application) AfterStop(fns ...func() error) {
	app.initialize()
	if app.afterStop == nil {
		app.afterStop = xdefer.NewStack()
	}
	app.afterStop.Push(fns...)
}

func (app *Application) Serve(s server.Server) error {
	app.servers = append(app.servers, s)
	return nil
}

func (app *Application) Schedule(w worker.Worker) error {
	app.workers = append(app.workers, w)
	return nil
}

// Job ..
func (app *Application) Job(runner job.Runner) error {
	namedJob, ok := runner.(interface{ GetJobName() string })
	// job runner must implement GetJobName
	if !ok {
		return nil
	}
	jobName := namedJob.GetJobName()
	if flag.Bool("disable-job") {
		xlog.Info("jupiter disable job", xlog.String("name", jobName))
		return nil
	}

	// start job by name
	jobFlag := flag.String("job")
	if jobFlag == "" {
		xlog.Error("jupiter jobs flag name empty", xlog.String("name", jobName))
		return nil
	}

	if jobName != jobFlag {
		xlog.Info("jupiter disable jobs", xlog.String("name", jobName))
		return nil
	}
	xlog.Info("jupiter register job", xlog.String("name", jobName))
	app.jobs[jobName] = runner
	return nil
}

func (app *Application) GracefulStop(ctx context.Context) (err error) {
	app.stopOnce.Do(func() {
		app.beforeStop.Clean()

		for _, s := range app.servers {
			func(s server.Server) {
				app.cycle.Run(func() error {
					return s.GracefulStop(ctx)
				})
			}(s)
		}
		for _, w := range app.workers {
			func(w worker.Worker) {
				app.cycle.Run(w.Stop)
			}(w)
		}

		select {
		case <-app.cycle.Done():
			app.cycle.Close()
		}
	})
	return err
}

func (app *Application) startServers() error {
	var eg errgroup.Group
	for _, s := range app.servers {
		s := s
		eg.Go(func() (err error) {
			xlog.Info("start servers", xlog.String("Label", s.Info().Label()))
			defer xlog.Info("exit server", xlog.FieldErr(err), xlog.String("Label", s.Info().Label()))
			return s.Serve()
		})
	}
	return eg.Wait()
}

func (app *Application) startWorkers() error {
	var eg errgroup.Group
	// start multi workers
	for _, w := range app.workers {
		w := w
		eg.Go(func() error {
			return w.Run()
		})
	}
	return eg.Wait()
}

func (app *Application) startJobs() error {
	if len(app.jobs) == 0 {
		return nil
	}
	var jobs = make([]func(), 0)
	// warp jobs
	for name, runner := range app.jobs {
		jobs = append(jobs, func() {
			xlog.Info("job run begin", xlog.String("name", name))
			defer xlog.Info("job run end", xlog.String("name", name))
			// runner.Run panic 错误在更上层抛出
			runner.Run()
		})
	}
	xgo.Parallel(jobs...)()
	return nil
}

func (app *Application) waitSignals() {
	xlog.Info("init listen signal", xlog.String("mod", "app"), xlog.String("event", "init"))
	signals.Shutdown(func(grace bool) {
		if grace {
			_ = app.GracefulStop(context.TODO())
		} else {
			_ = app.Stop()
		}
	})
}

func (app *Application) clean() {
	for i := len(app.defers) - 1; i >= 0; i-- {
		fn := app.defers[i]
		if err := fn(); err != nil {
			xlog.Error("clean.defer", xlog.String("func", xstring.FunctionName(fn)))
		}
	}
}

func (app *Application) initConfig() error {
	err := config.Init()
	if err != nil {
		panic(err)
	}
	return nil
}

func (app *Application) initLogger() error {
	err := xlog.Build()
	if err != nil {
		panic(err)
	}
	return nil
}

func (app *Application) initFrame() error {
	appConfig := config.GetAppConfig()

	if maxProcs := appConfig.MaxProc; maxProcs != 0 {
		runtime.GOMAXPROCS(maxProcs)
	} else {
		if _, err := maxprocs.Set(); err != nil {
			xlog.Error("auto max procs", xlog.FieldErr(err))
		}
	}
	xlog.Info("auto max procs", xlog.Int64("procs", int64(runtime.GOMAXPROCS(-1))))

	app.name = appConfig.Name
	app.version = appConfig.Version
	app.env = appConfig.Env

	return nil
}

func (app *Application) Name() string {
	return app.name
}

func (app *Application) Version() string {
	return app.version
}

func (app *Application) Env() string {
	return app.env
}

func (app *Application) printBanner() error {
	fmt.Println("Starting application ...")
	return nil
}

func (app *Application) printInitEnd() error {
	xlog.Info("microservice info", xlog.String("name", app.name), xlog.String("version", app.version), xlog.String("env", app.env))
	return nil
}
