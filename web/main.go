package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	srv := &http.Server{
		Addr:    ":8080",
		Handler: app.routes(),
	}

	go func() {
		infoLog.Printf("Запуск сервера на %s\n", srv.Addr)
		err := srv.ListenAndServe()
		errorLog.Fatal(err)
	}()

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt)
	<-sigc

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}
