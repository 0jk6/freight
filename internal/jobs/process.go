package jobs

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/0jk6/freight/internal/models"
	"github.com/0jk6/freight/internal/orchestrator"
	amqp "github.com/rabbitmq/amqp091-go"
)

type ContainerInfo struct {
	image     string
	command   []string
	namespace string
}

func ProcessSubmissions() {
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
	}

	defer ch.Close()

	q, err := ch.QueueDeclare("jobs_queue", true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	//consume messages

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Println("Failed to register a consumer")
	}

	//use the channel to run a goroutine
	var loop chan struct{}

	go func() {
		for d := range msgs {
			//spin up a container
			submission := models.SubmissionRequest{}
			json.Unmarshal(d.Body, &submission)
			spinUpJob(submission)
		}
	}()

	<-loop
}

// spin up a job on k8s
func spinUpJob(submission models.SubmissionRequest) {
	log.Printf("spinning up job for %s", submission.JobID)

	// //decode base64 data
	// base64Code, err := base64.StdEncoding.DecodeString(submission.Code)
	// if err != nil {
	// 	log.Println("Error decoding code")
	// 	return
	// }
	// submission.Code = string(base64Code)

	// log.Println(submission.Code)
	namespace := "freight-ns"

	containerInfo := createContainerInfo(submission)
	orchestrator.RunJob(namespace, submission.JobID, containerInfo.image, submission.Code, containerInfo.command)
}

func createContainerInfo(submission models.SubmissionRequest) ContainerInfo {
	containerInfo := ContainerInfo{}

	if submission.Language == "py" {
		containerInfo.image = "python:3.9.7"
		containerInfo.command = []string{"python", "-c", submission.Code}
	} else if submission.Language == "c" {
		containerInfo.image = "gcc:12.4.0"
		containerInfo.command = []string{
			"/bin/sh",
			"-c",
			fmt.Sprintf("echo '%s' > main.c && gcc main.c -o a.out && ./a.out", submission.Code),
		}
	} else if submission.Language == "cpp" {
		containerInfo.image = "gcc:12.4.0"
		containerInfo.command = []string{
			"/bin/sh",
			"-c",
			fmt.Sprintf("echo '%s' > main.cpp && g++ main.cpp -o a.out && ./a.out", submission.Code),
		}
	} else if submission.Language == "js" {
		containerInfo.image = "node:22-alpine3.19"
		containerInfo.command = []string{"node", "-e", submission.Code}
	} else if submission.Language == "go" {
		containerInfo.image = "golang:1.22.0"
		containerInfo.command = []string{
			"/bin/sh",
			"-c",
			fmt.Sprintf("echo '%s' > main.go && go build main.go && chmod +x main && ./main", submission.Code)}
	} else {
		containerInfo.image = ""
	}

	return containerInfo
}
