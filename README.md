# A simple job scheduler for running async tasks/functions.

To get it running, you first need to define the task that you want to execute in the form of a function.

```go
func myTask(ctx context.Context) error {
    //code goes here
}
```
