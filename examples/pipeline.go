package main

import "github.com/slavaavr/go-async"

func main2() {
	t1 := async.Submit(func() (string, error) {
		println("t1")
		return "A", nil
	})

	t2 := async.Submit(func() (string, error) {
		v, err := t1.Await()
		if err != nil {
			return "", err
		}

		println("t2")
		return v + "B", nil
	})

	t3 := async.Submit(func() (string, error) {
		v, err := t2.Await()
		if err != nil {
			return "", err
		}

		println("t3")
		return v + "C", nil
	})

	v, err := t3.Await()
	if err != nil {
		panic(err)
	}

	println(v == "ABC")
}
