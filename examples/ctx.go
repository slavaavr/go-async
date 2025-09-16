package main

import (
	"context"
	"errors"
	"time"

	"github.com/slavaavr/go-async"
)

func main() {
	v, err := runHeavyTasks(context.Background())
	if err != nil {
		time.Sleep(3 * time.Second)
		return
	}

	println(v)
}

func runHeavyTasks(ctx context.Context) (string, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	t1 := async.Submit(func() (string, error) {
		println("task 1")
		return "A", nil
	})

	t2 := async.Submit(func() (string, error) {
		println("task 2")
		return "", errors.New("some error")
	})

	t3 := async.Submit(func() (string, error) {
		select {
		case <-ctx.Done():
			println("ctx closed in task 3")
			return "", ctx.Err()

		case <-time.After(1 * time.Second):
			println("timeout in task 3")
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
