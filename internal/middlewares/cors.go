package middlewares

import (
	"net/http"
)

// Logger middleware
type Cors struct {
	handler http.Handler
}

func (c *Cors) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Add CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle preflight request
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Handle actual request
	c.handler.ServeHTTP(w, r)
}

func NewCors(handler http.Handler) *Cors {
	return &Cors{handler: handler}
}
