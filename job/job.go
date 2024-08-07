package job

import (
	"context"
	"github.com/google/uuid"
)

type Task func(ctx context.Context) error

// TODO: add retrial count
// add ability to start after some given time
// add ability to start only after a given job finished
type Job struct {
	id            string
	task          Task
	respond       bool
	response      chan Response
	setOwnChannel bool
	ctx           context.Context
}

func NewJob(task Task, options ...Option) Job {
	job := Job{
		id:            uuid.New().String(),
		task:          task,
		respond:       false,
		response:      make(chan Response),
		setOwnChannel: false,
		ctx:           context.Background(),
	}

	for _, option := range options {
		option(&job)
	}

	return job
}

func (j Job) AwaitResponse() Response {
	return <-j.response
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

func (j Job) Respond() bool {
	return j.respond
}

func (j Job) ResponseChannel() chan Response {
	return j.response
}

func (j Job) SetOwnChannel() bool {
	return j.setOwnChannel
}

type Option func(j *Job)

func WithID(id string) Option {
	return func(j *Job) {
		j.id = id
	}
}

func WithResponseChannel(response chan Response) Option {
	return func(j *Job) {
		j.response = response
		j.setOwnChannel = true
	}
}

func DoRespond() Option {
	return func(j *Job) {
		j.respond = true
	}
}

func WithContext(ctx context.Context) Option {
	return func(j *Job) {
		j.ctx = ctx
	}
}
