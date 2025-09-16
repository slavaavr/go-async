package async

import (
	"fmt"
	"sync/atomic"
)

type Task[T any] struct {
	out    chan *result[T]
	result atomic.Pointer[result[T]]
}

type result[T any] struct {
	value T
	err   error
}

func newResult[T any](value T, err error) *result[T] {
	return &result[T]{value: value, err: err}
}

func Submit[T any](
	f func() (T, error),
) *Task[T] {
	out := make(chan *result[T], 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				var v T
				out <- newResult(v, fmt.Errorf("panic: %v", r))
			}

			close(out)
		}()

		v, err := f()
		out <- newResult(v, err)
	}()

	return &Task[T]{
		out:    out,
		result: atomic.Pointer[result[T]]{},
	}
}

func (s *Task[T]) Await() (T, error) {
	r, ok := <-s.out
	if !ok {
		for {
			r = s.result.Load()
			if r != nil {
				return r.value, r.err
			}
		}
	}

	s.result.Store(r)

	return r.value, r.err
}

type ActionTask struct {
	task *Task[struct{}]
}

func (s *ActionTask) Await() error {
	_, err := s.task.Await()
	return err
}

func SubmitAction(
	f func() error,
) *ActionTask {
	return &ActionTask{
		task: Submit(func() (struct{}, error) {
			return struct{}{}, f()
		}),
	}
}
