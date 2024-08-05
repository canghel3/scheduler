package job

import "context"

type Task func(ctx context.Context) error

type Job struct {
	ID       string
	Task     Task
	Response chan error
	Context  context.Context
}

type JobResponse struct {
	ID  string
	Err error
}

func NewJob(task Task) Job {
	return Job{
		ID:       "some ID",
		Task:     task,
		Response: make(chan error),
	}
}

type Option func(j *Job)

func JobID(id string) Option {
	return func(j *Job) {
		j.ID = id
	}
}
