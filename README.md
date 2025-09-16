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
import "github.com/slavaavr/go-async"

task := async.Submit(func() (string, error) {
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
import "github.com/slavaavr/go-async"

task := async.SubmitAction(func() error {
	println("doing work without a response...")
	return nil
})

if err := task.Await(); err != nil {
	panic(err)
}
```

### Notes
- To prevent `goroutine leaks`, use `context` with `cancel/timeout` whenever possible ([example](https://github.com/slavaavr/go-async/tree/main/examples/ctx.go))
- For more usage cases, see the [examples](https://github.com/slavaavr/go-async/tree/main/examples) folder

[ci-badge]:      https://github.com/slavaavr/go-async/actions/workflows/main.yaml/badge.svg
[ci-runs]:       https://github.com/slavaavr/go-async/actions