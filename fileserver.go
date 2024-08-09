package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/unrolled/render"
	"github.com/urfave/negroni/v3"
	"golang.org/x/net/http2"
	"log/slog"
	"net/http"
	"time"
)

var renderer = render.New(render.Options{})

func main() {
	router := httprouter.New()
	router.GET("/", HomeHandler)
	router.POST("/persons", PersonsPostHandler)
	router.GET("/persons/:id", GetPersonsHandler)

	middleware := negroni.New()
	middleware.UseFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		start := time.Now()
		next(rw, r)
		slog.Info("request %s need %v", r.URL.Path, time.Since(start))
	})
	middleware.UseHandler(router)

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

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func HomeHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	_, err := fmt.Fprintln(w, "Welcome to the home page")
	if err != nil {
		slog.Error("Error writing to response: %#v", err)
	}
}

func PersonsPostHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	_, err := fmt.Fprintln(w, "Post page")
	if err != nil {
		slog.Error("Error in PersonsPostHandler: %#v", err)
	}
}

func GetPersonsHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	person := &Person{
		Name: "tu hu con",
		Age:  40,
	}
	err := renderer.JSON(w, http.StatusOK, person)
	if err != nil {
		slog.Error("Error in PersonsPostHandler: %#v", err)
	}
}
