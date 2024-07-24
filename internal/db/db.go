package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	pool *pgxpool.Pool
	once sync.Once
)

func SetupConnectionPool() {
	var err error

	postgresHost := os.Getenv("POSTGRES_HOST")

	if postgresHost == "" {
		postgresHost = "localhost" // Fallback to localhost if the environment variable is not set
	}

	log.Printf("POSTGRES HOST: %s", postgresHost)

	user := "postgres"     // your database user
	password := "password" // your database password
	dbname := "freight"    // your database name
	port := "5432"         // your database port

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, password, postgresHost, port, dbname)

	once.Do(func() {
		pool, err = pgxpool.New(context.Background(), connStr)

		if err != nil {
			log.Fatal(err)
		}

		err = pool.Ping(context.Background())
		if err != nil {
			log.Fatal(err)
		}

		// Setup the primary table
		createTableQuery := `
			CREATE TABLE IF NOT EXISTS submissions (
				id SERIAL PRIMARY KEY,
				language VARCHAR(50) NOT NULL,
				code TEXT NOT NULL,
				job_id UUID NOT NULL,
				output TEXT
			);
		`

		log.Println("Setting up the database table")
		_, err = pool.Exec(context.Background(), createTableQuery)

		if err != nil {
			log.Fatal(err)
		}
	})
}

func GetConnectionPool() *pgxpool.Pool {
	return pool
}
