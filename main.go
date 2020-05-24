package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
)

var (
	logger = log.New(os.Stdout, "logger: ", log.Ldate | log.Ltime | log.Lshortfile)
)

func main() {
	port := os.Getenv("FINAL_SCENES_PORT")
	if port == "" {
		logger.Fatal("FINAL_SCENES_PORT environment variable missing.")
	}
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("./www"))
	mux.Handle("/", LogWrapper(fs))

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
			logger.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	log.Println("listening....")
	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		logger.Fatalf("HTTP server ListenAndServe: %v", err)
	}

	<-idleConnsClosed
}
func LogWrapper(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		// Our middleware logic goes here...
		next.ServeHTTP(w, r)
	})
}
