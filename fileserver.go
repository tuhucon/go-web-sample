package main

import (
	"fileserver/internal/controller"
	"fileserver/internal/repository"
	"github.com/julienschmidt/httprouter"
	"github.com/unrolled/render"
	"github.com/urfave/negroni/v3"
	"golang.org/x/net/http2"
	"log/slog"
	"net/http"
	"time"
)

func main() {
	router := httprouter.New()

	middleware := negroni.New()
	middleware.UseFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		start := time.Now()
		next(rw, r)
		slog.Info("request %s need %v", r.URL.Path, time.Since(start))
	})
	middleware.UseHandler(router)

	renderer := render.New()
	helloController := controller.NewHelloController(renderer)
	personController := controller.NewPersonController(renderer, &repository.PersonRepositoryImp{})

	// handle hello api
	router.GET("/hello", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		if err := helloController.Hello(w, r); err != nil {
			slog.Error(err.Error())
		}
	})

	router.GET("/time", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		if err := helloController.Time(w, r); err != nil {
			slog.Error(err.Error())
		}
	})

	// handle person api
	router.GET("/persons", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		if err := personController.GetAll(w, r); err != nil {
			slog.Error(err.Error())
		}
	})

	// Create server
	server := &http.Server{
		Addr:    ":8080",
		Handler: middleware,
	}
	//curl -v --http2-prior-knowledge http://localhost:8080/persons/45  : test http2
	err := http2.ConfigureServer(server, &http2.Server{
		MaxHandlers:                  0,
		MaxConcurrentStreams:         0,
		MaxDecoderHeaderTableSize:    0,
		MaxEncoderHeaderTableSize:    0,
		MaxReadFrameSize:             0,
		PermitProhibitedCipherSuites: false,
		IdleTimeout:                  0,
		MaxUploadBufferPerConnection: 0,
		MaxUploadBufferPerStream:     0,
		NewWriteScheduler:            nil,
		CountError:                   nil,
	})
	if err != nil {
		slog.Error("Failed to configure http2 server: %v", err)
	}
	err = server.ListenAndServe()
	if err != nil {
		slog.Error("server is stopped: %#v", err)
	}
}
