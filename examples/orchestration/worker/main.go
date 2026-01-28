package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sijoma/camunda-go-sdk/orchestration"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	zeebeAddress := os.Getenv("ZEEBE_REST_ADDRESS")
	if zeebeAddress == "" {
		zeebeAddress = "http://localhost:8080"
	}

	baseURL, err := url.Parse(zeebeAddress)
	if err != nil {
		log.Fatalf("failed to parse zeebe address: %v", err)
	}

	c8, err := orchestration.NewClient(
		orchestration.WithBaseURL(*baseURL),
	)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}

	// Create a job worker
	worker := c8.NewJobWorker("my-task-type", handleJob,
		orchestration.WithWorkerName("example-worker"),
		orchestration.WithPollInterval(1*time.Second),
	)

	fmt.Println("Starting job worker for 'my-task-type'...")

	// Run the worker
	if err := worker.Run(ctx); err != nil && err != context.Canceled {
		log.Fatalf("worker error: %v", err)
	}

	fmt.Println("Worker stopped.")
}

func handleJob(client *orchestration.Client, job orchestration.ActivatedJob) {
	fmt.Printf("Handling job %s of type %s\n", job.Key, job.Type)

	// Do some work...
	fmt.Printf("Variables: %v\n", job.Variables)

	// Complete the job
	err := client.Job.Complete(context.Background(), job.Key, orchestration.JobCompleteRequest{
		Variables: map[string]any{
			"result": "processed",
		},
	})
	if err != nil {
		fmt.Printf("Failed to complete job %s: %v\n", job.Key, err)
	} else {
		fmt.Printf("Successfully completed job %s\n", job.Key)
	}
}
