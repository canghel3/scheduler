package queue

import (
	"github.com/Ginger955/scheduler/customerrors"
	"github.com/Ginger955/scheduler/job"
	"sync/atomic"
	"time"
)

type Queue struct {
	workers      int
	jobChannel   chan job.Job
	jobRegistry  map[string][]chan job.Job
	runningTasks atomic.Int32
}

func NewQueue(workers int) *Queue {
	q := &Queue{
		jobChannel: make(chan job.Job, workers),
		workers:    workers,
	}

	for i := 0; i < workers; i++ {
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

	q.jobChannel <- job
	q.runningTasks.Add(1)
}

func (q *Queue) processor() {
	for j := range q.jobChannel {
		func() {
			defer recovery(q, j)

			select {
			case <-j.Context().Done():
				//user cancelled or task timed out waiting in the queue
				//do nothing because response is handled in job.AwaitResponse
			default:
				data, err := j.Task()(j.Context())
				go respond(j, err, data)
			}

			q.runningTasks.Store(q.runningTasks.Load() - 1)
		}()
	}
}

func recovery(q *Queue, j job.Job) {
	if r := recover(); r != nil {
		//start a new processor since the previous one panicked and died.
		go q.processor()

		q.runningTasks.Store(q.runningTasks.Load() - 1)
		j.ResponseChannel() <- job.NewResponse(j.ID(), customerrors.NewRecoveredPanicError(r), nil)
	}
}

func respond(j job.Job, err error, data any) {
	j.ResponseChannel() <- job.NewResponse(j.ID(), err, data)
}
