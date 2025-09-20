package async

import (
	"fmt"
)

type Task[T any] struct {
	done   chan struct{}
	result T
	err    error
}

func Submit[T any](
	g *Group,
	f func() (T, error),
) *Task[T] {
	t := &Task[T]{done: make(chan struct{})}

	g.add()

	go func() {
		defer func() {
			if r := recover(); r != nil {
				t.err = fmt.Errorf("panic: %v", r)
			}

			close(t.done)
			g.done()
		}()

		t.result, t.err = f()
	}()

	return t
}

func (s *Task[T]) Await() (T, error) {
	<-s.done
	return s.result, s.err
}

type ActionTask struct {
	task *Task[struct{}]
}

func (s *ActionTask) Await() error {
	_, err := s.task.Await()
	return err
}

func SubmitAction(
	g *Group,
	f func() error,
) *ActionTask {
	return &ActionTask{
		task: Submit(g, func() (struct{}, error) {
			return struct{}{}, f()
		}),
	}
}
