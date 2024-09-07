package route

import (
	"esfgit.leju.com/golang/frame/server/xmux"
	"esfgit.leju.com/golang/frame/util/xswagger"
	_ "esfgit.leju.com/golang/frame/util/xswagger/statik"
	"github.com/rakyll/statik/fs"
	"net/http"
)

func RegisterSwagger(s *xmux.Server) {
	if s.Info().Scheme == "" {
		panic("ServiceInfo Scheme" + "不能为空")
	}
	xswagger.Services.AddHost(s.Info().Label(), s.Info().Scheme)
	//s.PathPrefix("/s/").Handler(http.StripPrefix("/s/", fs))
	s.HandleFunc("/q/services", func(w http.ResponseWriter, r *http.Request) {
		services, err := xswagger.Services.GetSwagger()
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(200)
		w.Write(services)
	}).Methods("GET")

	statikFS, err := fs.New()
	if err != nil {
		panic(err)
	}
	staticServer := http.FileServer(statikFS)
	sh := http.StripPrefix("/q/swagger-ui", staticServer)
	s.PathPrefix("/q/swagger-ui").Handler(sh)
}
