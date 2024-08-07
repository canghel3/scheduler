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

// Add adds a job to the queue.
func (q *Queue) Add(job job.Job) {
	//run add as routine because if the job channel is full, sending on it is blocked until it is read from, and we do not want to block the user
	go q.add(job)
}

func (q *Queue) add(job job.Job) {
	q.jobs <- job
}

func (q *Queue) processor() {
	for job := range q.jobs {
		func() {
			defer recovery(job)

			err := job.Task()(job.Context())
			go respond(job, err)
		}()
	}
}

func recovery(j job.Job) {
	if r := recover(); r != nil {
		if j.Respond() {
			j.ResponseChannel() <- job.NewResponse(j.ID(), customerrors.NewRecoveredPanicError(r))

			if !j.SetOwnChannel() {
				//job initiator is expecting a response, but has not set its own channel,
				//so it is safe to close the channel after sending the response.
				//the channel will be closed after the response is read.
				//in case the job initiator set its own response channel, the channel is not closed
				//since it may be used for reading other responses.
				close(j.ResponseChannel())
			}
		} else {
			//job initiator is not expecting a response, so close the channel
			close(j.ResponseChannel())
		}
	}
}

func respond(j job.Job, err error) {
	if j.Respond() {
		j.ResponseChannel() <- job.NewResponse(j.ID(), err)

		if !j.SetOwnChannel() {
			//job initiator is expecting a response, but has not set its own channel,
			//so it is safe to close the channel after sending the response.
			//the channel will be closed after the response is read.
			//in case the job initiator set its own response channel, the channel is not closed
			//since it may be used for reading other responses.
			close(j.ResponseChannel())
		}
	} else {
		//job initiator is not expecting a response, so close the channel
		close(j.ResponseChannel())
	}
}
