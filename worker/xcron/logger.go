package xcron

import "github.com/Remember9/frame/xlog"

var logMod = xlog.String("mod", "xcron")

type wrappedLogger struct {
}

// Info logs routine messages about cron's operation.
func (wl *wrappedLogger) Info(msg string, keysAndValues ...interface{}) {
	keysAndValues = append(keysAndValues, logMod)
	xlog.Infow("cron "+msg, keysAndValues...)
}

// Error logs an error condition.
func (wl *wrappedLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	keysAndValues = append(keysAndValues, logMod, xlog.FieldErr(err))
	xlog.Errorw("cron2 "+msg, keysAndValues...)
}
