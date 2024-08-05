package queue

import (
	"github.com/Ginger955/scheduler/customerrors"
	"github.com/Ginger955/scheduler/job"
)

type Queue struct {
	jobs chan job.Job
	size int
}

func NewQueue(size int) *Queue {
	return &Queue{
		jobs: make(chan job.Job, size),
		size: size,
	}
}

func (q *Queue) Start() {
	for i := 0; i < q.size; i++ {
		go q.processor()
	}
}

func (q *Queue) Add(job job.Job) {
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
		job.Response <- customerrors.NewRecoveredPanic(r)
	}
}

func respond(job job.Job, err error) {
	job.Response <- err
}
