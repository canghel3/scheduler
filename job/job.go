package job

import (
	"context"
	"github.com/google/uuid"
)

type Task func(ctx context.Context) error

type Job struct {
	id       string
	task     Task
	response chan Response
	ctx      context.Context
}

func NewJob(task Task, options ...Option) Job {
	job := Job{
		id:       uuid.New().String(),
		task:     task,
		response: make(chan Response),
		ctx:      context.Background(),
	}

	for _, option := range options {
		option(&job)
	}

	return job
}

func (j Job) AwaitResponse() Response {
	s := <-j.response
	//the receiver should never close the channel, but in this case the channel is single use only, and it will never be used again for sending
	close(j.response)
	return s
}

func (j Job) ID() string {
	return j.id
}

func (j Job) Task() Task {
	return j.task
}

func (j Job) Context() context.Context {
	return j.ctx
}

func (j Job) ResponseChannel() chan Response {
	return j.response
}

type Option func(j *Job)

func WithID(id string) Option {
	return func(j *Job) {
		j.id = id
	}
}

func WithContext(ctx context.Context) Option {
	return func(j *Job) {
		j.ctx = ctx
	}
}
