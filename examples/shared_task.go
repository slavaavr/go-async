package main

import "github.com/slavaavr/go-async"

func main3() {
	sharedTask := async.Submit(func() (string, error) {
		println("shared task")
		return "A", nil
	})

	t1 := async.Submit(func() (string, error) {
		v, err := sharedTask.Await()
		if err != nil {
			return "", err
		}

		println("t1")
		return v + "B", nil
	})

	t2 := async.Submit(func() (string, error) {
		v, err := sharedTask.Await()
		if err != nil {
			return "", err
		}

		println("t2")
		return v + "C", nil
	})

	t3 := async.Submit(func() (string, error) {
		v1, err := t1.Await()
		if err != nil {
			return "", err
		}

		v2, err := t2.Await()
		if err != nil {
			return "", err
		}

		println("t3")
		return v1 + v2, nil
	})

	v, err := t3.Await()
	if err != nil {
		panic(err)
	}

	println(v == "ABAC")
}
