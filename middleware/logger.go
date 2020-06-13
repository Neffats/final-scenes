package middleware

import (
	"log"
	"net/http"
)

type LogWrapper struct {
	Logger *log.Logger
}

func (lw *LogWrapper) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lw.Logger.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		// Our middleware logic goes here...
		next.ServeHTTP(w, r)
	})
}
