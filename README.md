## A simple Go job scheduler for running async tasks/functions.

### Install with

```bash
go get github.com/Ginger955/scheduler
```

### How to use:

1. define the task that you want to execute in the form of a function.

```go
package main

import "context"

func myTask(ctx context.Context) (data any, err error) {
	//code goes here
	return "data", nil
}
```

2. initialize a queue

```go
package main

import (
	"github.com/Ginger955/scheduler/queue"
	"context"
)

func myTask(ctx context.Context) (data any, err error) {
	//code goes here
	return "data", nil
}

func main() {
	//a worker refers to a routine that processes a task
	//queue with 3 workers
	q := queue.NewQueue(3)
}
```

3. initialize a job and add it to the queue

```go
package main

import (
	"github.com/Ginger955/scheduler/job"
	"github.com/Ginger955/scheduler/queue"
	"context"
)

func myTask(ctx context.Context) (data any, err error) {
	//code goes here
	return "data", nil
}

func main() {
	//a worker refers to a routine that processes a task
	//queue with 3 workers
	q := queue.NewQueue(3)

	job1 := job.NewJob(myTask)
	//add job1 to the queue
	q.Add(job1)

	//you can also pass a context to a job
	ctx, cancel := context.WithCancel(context.WithValue(context.Background(), "data", "some data you want to pass"))
	defer cancel()

	//passing the context to the job
	job2 := job.NewJob(myTask, job.WithContext(ctx))
	q.Add(job2)
}
```

Once a job has been queued, it will be processed as soon as a worker is free and receives the job. \
The result of the job execution can be checked as such:

```go
package main

import (
	"fmt"
	"github.com/Ginger955/scheduler/job"
	"github.com/Ginger955/scheduler/queue"
	"context"
)

func myTask(ctx context.Context) (data any, err error) {
	//code goes here
	return "data", nil
}

func main() {
	//a worker refers to a routine that processes a task
	//queue with 3 workers
	q := queue.NewQueue(3)

	job1 := job.NewJob(myTask)
	//add job1 to the queue
	q.Add(job1)

	//you can also pass a context to a job
	ctx, cancel := context.WithCancel(context.WithValue(context.Background(), "data", "some data you want to pass"))
	defer cancel()

	//passing the context to the job
	job2 := job.NewJob(myTask, job.WithContext(ctx))
	q.Add(job2)

	//note that this operation is blocking
	response := job2.AwaitResponse()

	if response.Err() != nil {
		//handle error
	}

	fmt.Println(response.ID())
	fmt.Println(response.Data())
}
```