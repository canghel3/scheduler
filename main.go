package main

import (
	"context"
	"fmt"
	"github.com/Ginger955/scheduler/job"
	"github.com/Ginger955/scheduler/queue"
	"time"
)

func main() {
	q := queue.NewQueue(5)
	q.Start()

	ctx, cancel := context.WithCancel(context.WithValue(context.Background(), "number", "this is not a number"))
	//cancel()
	defer cancel()

	j1 := job.NewJob(task1, job.WithContext(ctx))
	j2 := job.NewJob(task2)
	j3 := job.NewJob(samplePassValue, job.WithContext(ctx))

	q.Add(j1)
	q.Add(j2)
	q.Add(j3)

	r := j1.AwaitResponse()
	fmt.Println(r.ID())

	err := r.Err()
	if err != nil {
		fmt.Println(err.Error())
	}

	time.Sleep(time.Second * 1)
}

func task1(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:

	}
	panic("implement me")
	fmt.Println("task 1")
	return nil
}

func task2(context.Context) error {
	fmt.Println("task 2")
	return nil
}

func samplePassValue(ctx context.Context) error {
	fmt.Println(ctx.Value("number"))
	return nil
}
