package middlewares

import (
	"net/http"
)

// Logger middleware
type Cors struct {
	handler http.Handler
}

func (c *Cors) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//add cors headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	c.handler.ServeHTTP(w, r)
}

func NewCors(handler http.Handler) *Cors {
	return &Cors{handler: handler}
}
