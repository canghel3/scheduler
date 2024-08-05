package main

import (
	"context"
	"fmt"
	"github.com/Ginger955/scheduler/job"
	"github.com/Ginger955/scheduler/queue"
	"time"
)

func main() {
	q := queue.NewQueue(1)
	q.Start()

	j1 := job.NewJob(task1)
	j2 := job.NewJob(task2)

	ctx, cancel := context.WithCancel(context.Background())
	j1.Context = ctx
	cancel()
	//defer cancel()

	q.Add(j1)
	q.Add(j2)

	//TODO: if any response is sent on the response channel, but is not read, the routine is blocked.
	// this can cause high resource utilisation if many are left unread.
	//j1err := <-j1.Response
	//fmt.Println(j1err.Error())

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
