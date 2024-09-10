package queue

import (
	"github.com/Ginger955/scheduler/customerrors"
	"github.com/Ginger955/scheduler/job"
	"time"
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

// Add adds a job to the queue.
func (q *Queue) Add(job job.Job, delay ...time.Duration) {
	//run add as routine because if the job channel is full, sending on it is blocked until it is read from, and we do not want to block the user
	go q.add(job, delay...)
}

func (q *Queue) add(job job.Job, delay ...time.Duration) {
	if len(delay) > 0 {
		time.Sleep(delay[0])
	}
	q.jobs <- job
}

func (q *Queue) processor() {
	for j := range q.jobs {
		func() {
			defer recovery(j)

			err := j.Task()(j.Context())
			go respond(j, err)
		}()
	}
}

func recovery(j job.Job) {
	if r := recover(); r != nil {
		j.ResponseChannel() <- job.NewResponse(j.ID(), customerrors.NewRecoveredPanicError(r))
	}
}

func respond(j job.Job, err error) {
	j.ResponseChannel() <- job.NewResponse(j.ID(), err)
}
