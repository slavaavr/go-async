package main

import (
	"context"
	"errors"
	"time"

	"github.com/slavaavr/go-async"
)

func main() {
	v, err := runTasks(context.Background())
	if err != nil {
		println("got expected error:", err.Error())
		return
	}

	println(v)
}

func runTasks(ctx context.Context) (string, error) {
	group, ctx := async.NewGroup(ctx)
	defer group.Close()

	t1 := async.Submit(group, func() (string, error) {
		println("task 1")
		return "A", nil
	})

	t2 := async.Submit(group, func() (string, error) {
		println("task 2")
		return "", errors.New("some error")
	})

	t3 := async.Submit(group, func() (string, error) {
		select {
		case <-ctx.Done():
			println("ctx closed in task 3")
			return "", ctx.Err()

		case <-time.After(1 * time.Second):
			println("task 3")
		}

		return "C", nil
	})

	v1, err := t1.Await()
	if err != nil {
		return "", err
	}

	v2, err := t2.Await()
	if err != nil {
		return "", err
	}

	v3, err := t3.Await()
	if err != nil {
		return "", err
	}

	return v1 + v2 + v3, nil
}
