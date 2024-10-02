package queue

import (
	"github.com/Ginger955/scheduler/customerrors"
	"github.com/Ginger955/scheduler/job"
	"sync/atomic"
	"time"
)

type Queue struct {
	jobs         chan job.Job
	runningTasks atomic.Int32
	workers      int
}

func NewQueue(size int, workers int) *Queue {
	q := &Queue{
		jobs:    make(chan job.Job, size),
		workers: workers,
	}

	for i := 0; i < q.workers; i++ {
		go q.processor()
	}

	return q
}

// Add adds a job to the queue.
func (q *Queue) Add(job job.Job, delay ...time.Duration) {
	//run add as routine because if the job channel is full, sending on it is blocked until it is read from, and we do not want to block the user
	go q.add(job, delay...)
}

func (q *Queue) Running() int {
	return int(q.runningTasks.Load())
}

func (q *Queue) add(job job.Job, delay ...time.Duration) {
	if len(delay) > 0 {
		time.Sleep(delay[0])
	}

	q.jobs <- job
	q.runningTasks.Add(1)
}

func (q *Queue) processor() {
	for j := range q.jobs {
		func() {
			defer recovery(q, j)

			err := j.Task()(j.Context())
			go respond(j, err)
			q.runningTasks.Store(q.runningTasks.Load() - 1)
		}()
	}
}

func recovery(q *Queue, j job.Job) {
	if r := recover(); r != nil {
		q.runningTasks.Store(q.runningTasks.Load() - 1)
		j.ResponseChannel() <- job.NewResponse(j.ID(), customerrors.NewRecoveredPanicError(r))
	}
}

func respond(j job.Job, err error) {
	j.ResponseChannel() <- job.NewResponse(j.ID(), err)
}
