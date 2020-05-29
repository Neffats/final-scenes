package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/Neffats/final-scenes/handlers"
	"github.com/Neffats/final-scenes/middleware"
	"github.com/Neffats/final-scenes/stores"
)

var (
	logger = log.New(os.Stdout, "logger: ", log.Ldate|log.Ltime|log.Lshortfile)
)

func main() {
	store := stores.NewFilmStore("films.json")
	err := store.Init()
	if err != nil {
		logger.Fatalf("failed to initialise film store: %v", err)
	}

	h := &handlers.HTTP{
		Films:  store,
		Logger: logger,
	}
	l := &middleware.LogWrapper{
		Logger: logger,
	}

	port := os.Getenv("FINAL_SCENES_PORT")
	if port == "" {
		logger.Fatal("FINAL_SCENES_PORT environment variable missing.")
	}
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("./www"))
	mux.Handle("/static/", l.Wrap(http.StripPrefix("/static/", fs)))
	mux.Handle("/guess/", l.Wrap(http.HandlerFunc(h.HandleGuess)))
	mux.Handle("/", l.Wrap(http.HandlerFunc(h.HandleTemplate)))

	srv := &http.Server{
		Addr:     fmt.Sprintf("0.0.0.0:%s", port),
		Handler:  mux,
		ErrorLog: logger,
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		if err := srv.Shutdown(context.Background()); err != nil {
			logger.Printf("HTTP server Shutdown: %v\n", err)
		}
		close(idleConnsClosed)
	}()

	logger.Println("listening....")
	if err = srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		logger.Fatalf("HTTP server ListenAndServe: %v\n", err)
	}

	<-idleConnsClosed
}
