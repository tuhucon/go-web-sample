package main

import (
	"database/sql"
	"errors"
	"fileserver/internal/controller"
	"fileserver/internal/repository"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
	"github.com/unrolled/render"
	"github.com/urfave/negroni/v3"
	_ "go.uber.org/automaxprocs"
	"golang.org/x/net/http2"
)

func main() {
	// mysql DB
	var db *sql.DB
	var err error
	// Capture connection properties.
	cfg := mysql.Config{
		User:   "root",
		Passwd: "tuhucon",
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "test",
	}
	if db, err = sql.Open("mysql", cfg.FormatDSN()); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()
	// seting connection pool, max lifetime
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(10)
	db.SetConnMaxLifetime(8 * time.Hour)
	rows, err := db.Query("select * from person")
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() { // if for end normally or break by an error, rows.close() is call automatically
		var id int
		var name string
		var age int
		rows.Scan(&id, &name, &age)
		fmt.Println(id, name, age)
	}
	rows.Close() // we can call rows.Close() multiple time, ignoring its err is best practices.
	//rows.error() store the error raised in for loop, so always check rows.error() after loop and ignore err in for loop
	if err = rows.Err(); err != nil {
		slog.Error(err.Error())
	}

	router := httprouter.New()

	middleware := negroni.New()
	middleware.UseFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		start := time.Now()
		next(rw, r)
		slog.Info("request %s need %v", r.URL.Path, time.Since(start))
	})

	var name string
	err = db.QueryRow("select name from person where id = ?", 1).Scan(&name) // QueryRow defers error untils we call Scan
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) { // ErrNoRows isn't error, we should treat it as a biz logic use-case.
			// there were no rows, but otherwise no error occurred
		} else {
			log.Fatal(err)
		}
	}

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
