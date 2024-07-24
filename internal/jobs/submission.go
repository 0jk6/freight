package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/0jk6/freight/internal/db"
	"github.com/0jk6/freight/internal/models"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

// submit a job and store it in a queue and return the job id
func Push(submission models.SubmissionRequest) string {
	//connect to rabbitmq and push the submission to the queue
	rabbitHost := os.Getenv("RABBITMQ_HOST")

	if rabbitHost == "" {
		rabbitHost = "localhost" // Fallback to localhost if the environment variable is not set
	}

	rabbitUser := "guest"
	rabbitPassword := "guest"
	rabbitPort := "5672"

	rabbitConnStr := fmt.Sprintf("amqp://%s:%s@%s:%s/", rabbitUser, rabbitPassword, rabbitHost, rabbitPort)

	conn, err := amqp.Dial(rabbitConnStr)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Println("Error while creating channel.")
		return ""
	}
	defer ch.Close()

	//publish the submission to the queue
	q, err := ch.QueueDeclare("jobs_queue", true, false, false, false, nil)
	if err != nil {
		log.Println("Error while declaring queue.")
		log.Println(err)
		return ""
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//generate a unique id for the job
	submission.JobID = generateUUID()

	msg, err := json.Marshal(submission)
	if err != nil {
		log.Println("Cannot encode message")
		return ""
	}

	err = ch.PublishWithContext(ctx, "", q.Name, false, false, amqp.Publishing{ContentType: "text/plain", Body: msg})

	if err != nil {
		log.Println("Error while publishing message.")
		return ""
	}

	//store it in the db
	err = storeJob(submission)
	if err != nil {
		log.Println(err)
		return submission.JobID
	}

	return submission.JobID
}

func generateUUID() string {
	return uuid.New().String()
}

func storeJob(submission models.SubmissionRequest) error {
	pool := db.GetConnectionPool()
	_, err := pool.Exec(context.Background(), "INSERT INTO submissions (language, code, job_id, output) VALUES ($1, $2, $3, $4)", submission.Language, submission.Code, submission.JobID, "")

	return err
}
