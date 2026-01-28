package orchestration

import (
	"context"
	"time"
)

type JobHandler func(client *Client, job ActivatedJob)

type Worker struct {
	client            *Client
	jobType           string
	handler           JobHandler
	workerName        string
	timeout           int64
	maxJobsToActivate int32
	pollInterval      time.Duration
	fetchVariables    []string
}

type WorkerOption func(*Worker)

func WithWorkerName(name string) WorkerOption {
	return func(w *Worker) {
		w.workerName = name
	}
}

func WithTimeout(timeout time.Duration) WorkerOption {
	return func(w *Worker) {
		w.timeout = int64(timeout.Milliseconds())
	}
}

func WithMaxJobsToActivate(n int32) WorkerOption {
	return func(w *Worker) {
		w.maxJobsToActivate = n
	}
}

func WithPollInterval(interval time.Duration) WorkerOption {
	return func(w *Worker) {
		w.pollInterval = interval
	}
}

func WithFetchVariables(vars ...string) WorkerOption {
	return func(w *Worker) {
		w.fetchVariables = vars
	}
}

func (c *Client) NewJobWorker(jobType string, handler JobHandler, opts ...WorkerOption) *Worker {
	w := &Worker{
		client:            c,
		jobType:           jobType,
		handler:           handler,
		workerName:        "default",
		timeout:           30000,
		maxJobsToActivate: 1,
		pollInterval:      100 * time.Millisecond,
	}

	for _, opt := range opts {
		opt(w)
	}

	return w
}

func (w *Worker) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			jobs, err := w.client.Job.Activate(ctx, JobActivateRequest{
				Type:              w.jobType,
				Worker:            w.workerName,
				Timeout:           w.timeout,
				MaxJobsToActivate: w.maxJobsToActivate,
				FetchVariable:     w.fetchVariables,
			})
			if err != nil {
				// In case of error, wait poll interval before retrying
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(w.pollInterval):
					continue
				}
			}

			if len(jobs) == 0 {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(w.pollInterval):
					continue
				}
			}

			for _, job := range jobs {
				w.handler(w.client, job)
			}
		}
	}
}
