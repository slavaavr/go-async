// nolint:unused
package main

import (
	"context"
	"fmt"

	"github.com/slavaavr/go-async"
)

func main3() {
	group, ctx := async.NewGroup(context.Background())
	_ = ctx
	defer group.Close()

	sharedTask := async.Submit(group, func() (string, error) {
		return "A", nil
	})

	t1 := async.Submit(group, func() (string, error) {
		v, err := sharedTask.Await()
		if err != nil {
			return "", err
		}

		return v + "B", nil
	})

	t2 := async.Submit(group, func() (string, error) {
		v, err := sharedTask.Await()
		if err != nil {
			return "", err
		}

		return v + "C", nil
	})

	t3 := async.Submit(group, func() (string, error) {
		v1, err := t1.Await()
		if err != nil {
			return "", err
		}

		v2, err := t2.Await()
		if err != nil {
			return "", err
		}

		return v1 + v2, nil
	})

	res, err := t3.Await()
	if err != nil {
		panic(err)
	}

	fmt.Printf("expect='%s', actual='%s'\n", "ABAC", res)
}
