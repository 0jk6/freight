package main

import (
	"log"
	"net/http"

	"github.com/0jk6/freight/internal/db"
	"github.com/0jk6/freight/internal/handlers"
	"github.com/0jk6/freight/internal/middlewares"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", handlers.HomeHandler)
	mux.HandleFunc("POST /submission", handlers.SubmissionHandler)
	mux.HandleFunc("POST /output", handlers.OutputHandler)
	mux.HandleFunc("GET /check", handlers.CheckHandler)
	wrappedMux := middlewares.NewLogger(middlewares.NewCors(mux))

	//setup the database connection pool
	log.Println("Setting up database connection pool")
	db.SetupConnectionPool()

	log.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", wrappedMux))
}
