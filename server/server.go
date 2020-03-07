package server

import (
	"fmt"

	"github.com/NYTimes/gziphandler"

	"github.com/AsCat/acorn/config"
	"github.com/AsCat/acorn/log"
	"github.com/AsCat/acorn/routing"
	"net/http"
)

type Server struct {
	httpServer *http.Server
}

func NewServer() *Server {
	conf := config.Get()

	router := routing.NewRouter()

	if conf.Server.CORSAllowAll {
		router.Use(corsAllowed)
	}

	handler := http.Handler(router)
	if conf.Server.GzipEnabled {
		handler = configureGzipHandler(router)
	}

	// The Kiali server has only a single http server ever during its lifetime. But to support
	// testing that wants to start multiple servers over the lifetime of the process,
	// we need to override the default server mux with a new one everytime.
	mux := http.NewServeMux()
	http.DefaultServeMux = mux
	http.Handle("/", handler)

	// create the server definition that will handle both console and api server traffic
	httpServer := &http.Server{
		Addr: fmt.Sprintf("%v:%v", conf.Server.Address, conf.Server.Port),
	}

	// return our new Server
	return &Server{
		httpServer: httpServer,
	}
}

func (s *Server) Start() {
	conf := config.Get()
	log.Infof("Server endpoint will start at [%v%v]", s.httpServer.Addr, conf.Server.WebRoot)
	log.Infof("Server endpoint will serve static content from [%v]", conf.Server.StaticContentRootDirectory)
	secure := conf.Identity.CertFile != "" && conf.Identity.PrivateKeyFile != ""
	go func() {
		var err error
		if secure {
			log.Infof("Server endpoint will require https")
			err = s.httpServer.ListenAndServeTLS(conf.Identity.CertFile, conf.Identity.PrivateKeyFile)
		} else {
			err = s.httpServer.ListenAndServe()
		}
		log.Warning(err)
	}()

	// Start the Metrics Server
	//if conf.Server.MetricsEnabled {
	//	StartMetricsServer()
	//}
}

// Stop the HTTP server
func (s *Server) Stop() {
	//business.Stop()
	log.Infof("Server endpoint will stop at [%v]", s.httpServer.Addr)
	s.httpServer.Close()
}

func corsAllowed(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
		next.ServeHTTP(w, r)
	})
}

func configureGzipHandler(handler http.Handler) http.Handler {
	contentTypeOption := gziphandler.ContentTypes([]string{
		"application/javascript",
		"application/json",
		"image/svg+xml",
		"text/css",
		"text/html",
	})
	if handlerFunc, err := gziphandler.GzipHandlerWithOpts(contentTypeOption); err == nil {
		return handlerFunc(handler)
	} else {
		// This could happen by a wrong configuration being sent to GzipHandlerWithOpts
		panic(err)
	}
}
