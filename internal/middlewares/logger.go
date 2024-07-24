package middlewares

import (
	"log"
	"net/http"
	"time"
)

// Logger middleware
type Logger struct {
	handler http.Handler
}

func (l *Logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	l.handler.ServeHTTP(w, r)
	log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(startTime).String())
}

func NewLogger(handler http.Handler) *Logger {
	return &Logger{handler: handler}
}
