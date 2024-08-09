package main

import (
	"database/sql"
	"fileserver/internal/controller"
	"fileserver/internal/repository"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
	"github.com/unrolled/render"
	"github.com/urfave/negroni/v3"
	"golang.org/x/net/http2"
	"log/slog"
	"net/http"
	"time"
)

func main() {
	// mysql DB
	var db *sql.DB
	var err error
	if db, err = sql.Open("mysql", "root:tuhucon@tcp(127.0.0.1:3306)/test"); err != nil {
		slog.Error(err.Error())
	}
	defer db.Close()

	rows, err := db.Query("select * from person")
	defer rows.Close()
	if err == nil {
		for rows.Next() {
			var id int
			var name string
			var age int
			rows.Scan(&id, &name, &age)
			fmt.Println(id, name, age)
		}
	} else {
		slog.Error(err.Error())
	}

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
	err = http2.ConfigureServer(server, &http2.Server{
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
