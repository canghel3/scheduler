package job

import (
	"context"
	"fmt"
	"github.com/google/uuid"
)

type Task func(ctx context.Context) (data any, err error)

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

// AwaitResponse blocks until a response is received from the job execution.
// It is guaranteed to receive a response, even if the executing job panics.
func (j Job) AwaitResponse() Response {
	select {
	case <-j.ctx.Done():
		return Response{
			id:   j.id,
			err:  j.ctx.Err(),
			data: nil,
		}
	case response, ok := <-j.response:
		if !ok {
			return Response{
				id:   j.id,
				err:  fmt.Errorf("response channel closed before a response was received"),
				data: nil,
			}
		}
		return response
	}
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
