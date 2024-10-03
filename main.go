package main

import (
	"context"
	"fmt"
	"github.com/Ginger955/scheduler/job"
	"github.com/Ginger955/scheduler/queue"
	"time"
)

func main() {
	q := queue.NewQueue(2)

	ctx, cancel := context.WithCancel(context.WithValue(context.Background(), "number", "this is not a number"))
	//cancel()
	defer cancel()

	j1 := job.NewJob(task1, job.WithContext(ctx))
	j2 := job.NewJob(task2)
	j3 := job.NewJob(samplePassValue, job.WithContext(ctx))

	q.Add(j1)
	q.Add(j2)
	q.Add(j2)
	q.Add(j2)
	q.Add(j2)
	time.Sleep(100 * time.Millisecond)
	q.Add(j3)

	r := j1.AwaitResponse()
	fmt.Println(r.ID())

	err := r.Err()
	if err != nil {
		fmt.Println(err.Error())
	}

	r2 := j2.AwaitResponse()
	if r2.Err() != nil {
		fmt.Println(r2.Err().Error())
	}

	fmt.Println("r2data", r2.Data())

	r3 := j3.AwaitResponse()
	if r3.Err() != nil {
		fmt.Println(r3.Err())
	}
	fmt.Println(r3.Data())
}

func task1(ctx context.Context) (data any, err error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:

	}
	panic("panic 1")
	//fmt.Println("task 1")
	return nil, nil
}

func task2(context.Context) (data any, err error) {
	//fmt.Println("task 2")
	panic("panic 2")
	return "some random data", nil
}

func samplePassValue(ctx context.Context) (data any, err error) {
	time.Sleep(2 * time.Second)
	fmt.Println(ctx.Value("number"))
	return nil, nil
}
