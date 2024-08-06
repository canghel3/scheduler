package queue

import (
	"github.com/Ginger955/scheduler/customerrors"
	"github.com/Ginger955/scheduler/job"
)

type Queue struct {
	jobs    chan job.Job
	workers int
}

func NewQueue(size int) *Queue {
	return &Queue{
		jobs:    make(chan job.Job, size),
		workers: size,
	}
}

func (q *Queue) Start() {
	for i := 0; i < q.workers; i++ {
		go q.processor()
	}
}

func (q *Queue) Add(job job.Job) {
	//TODO: if queue is full, this operation is blocking (run as routine?)
	q.jobs <- job
}

func (q *Queue) processor() {
	for job := range q.jobs {
		func() {
			defer recovery(job)

			err := job.Task(job.Context)
			go respond(job, err)
		}()
	}
}

func recovery(job job.Job) {
	if r := recover(); r != nil {
		job.Response <- customerrors.NewRecoveredPanicError(r)
	}
}

func respond(job job.Job, err error) {
	job.Response <- err
}
