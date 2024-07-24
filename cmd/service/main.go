package main

import (
	"log"
	"time"

	"github.com/0jk6/freight/internal/db"
	"github.com/0jk6/freight/internal/jobs"
	"github.com/0jk6/freight/internal/orchestrator"
)

//this is another microservice that will run the jobs on kubernetes

func main() {
	//start processing the jobs in the queue
	log.Println("Starting job processor")

	db.SetupConnectionPool()

	namespace := "freight-ns"

	//check pod logs every 10 seconds
	go func() {
		for {
			orchestrator.ListJobs(namespace)
			time.Sleep(5 * time.Second)
		}
	}()

	jobs.ProcessSubmissions()

}
