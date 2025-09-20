package main

import (
	"context"
	"fmt"

	"github.com/slavaavr/go-async"
)

func main2() {
	group, ctx := async.NewGroup(context.Background())
	_ = ctx
	defer group.Close()

	t1 := async.Submit(group, func() (string, error) {
		return "A", nil
	})

	t2 := async.Submit(group, func() (string, error) {
		v, err := t1.Await()
		if err != nil {
			return "", err
		}

		return v + "B", nil
	})

	t3 := async.Submit(group, func() (string, error) {
		v, err := t2.Await()
		if err != nil {
			return "", err
		}

		return v + "C", nil
	})

	res, err := t3.Await()
	if err != nil {
		panic(err)
	}

	fmt.Printf("expect='%s', actual='%s'\n", "ABC", res)
}
