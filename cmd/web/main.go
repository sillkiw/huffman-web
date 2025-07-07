package main

import (
	"flag"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/sillkiw/huffman-web/internal/app"
)

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	app := app.NewApp(logger)

	srv := http.Server{
		Addr:         *addr,
		ErrorLog:     errorLog,
		Handler:      app.Routes(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	logger.Info("Run server", slog.String("addr", *addr))
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
