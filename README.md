## A simple job scheduler for running async tasks/functions.

How to run:

1. define the task that you want to execute in the form of a function.

```go
package main

import "context"

func myTask(ctx context.Context) error {
	//code goes here
	return nil
}
```

2. initialize a queue

```go
package main

import (
	"github.com/Ginger955/scheduler/queue"
	"context"
)

func myTask(ctx context.Context) error {
	//code goes here
	return nil
}

func main() {
	//a worker refers to a routine that processes a task
	//queue with size 2 (and 2 workers)
	q := queue.NewQueue(2)

	//queue with size 2 and 3 workers
	q = queue.NewQueue(2, 3)
}
```

3. 
