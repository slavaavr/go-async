## go-async
`go-async` is a library that aims to simplify the way of working with multiple concurrent tasks 


[![Build Status][ci-badge]][ci-runs]

### Installation

```sh
$ go get -u github.com/slavaavr/go-async
```

### Quick start
- Task with a signature `func() (T, error)`
```go
group, ctx := async.NewGroup(context.Background())
defer group.Close()

task := async.Submit(group, func() (string, error) {
    select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
    }
	return "hello World", nil
})

v, err := task.Await()
if err != nil {
	panic(err)
}

println(v)
```
- Action Task with a signature `func() error`
```go
group, ctx := async.NewGroup(context.Background())
defer group.Close()

task := async.SubmitAction(group, func() error {
    select {
    case <-ctx.Done():
        return ctx.Err()
    default:
    }
	println("executing a task without a response...")
	return nil
})

if err := task.Await(); err != nil {
	panic(err)
}
```

### Notes
- The `Group` entity ensures that all tasks are properly waited on, preventing `goroutine leaks` even if one task returns an `error` and stops further execution.
- For more use cases, see the [examples](https://github.com/slavaavr/go-async/tree/main/examples) folder.

[ci-badge]:      https://github.com/slavaavr/go-async/actions/workflows/main.yaml/badge.svg
[ci-runs]:       https://github.com/slavaavr/go-async/actions