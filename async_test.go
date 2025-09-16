package async

import (
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTask_Await(t *testing.T) {
	type customType struct {
		A int
		B string
		C bool
	}

	type response struct {
		A int
		B string
		C bool
		D time.Time
		E customType
	}

	res := response{
		A: 42,
		B: "some value",
		C: true,
		D: time.Date(2025, time.September, 13, 13, 13, 13, 13, time.UTC),
		E: customType{
			A: 1,
			B: "2",
			C: true,
		},
	}

	cases := []struct {
		name        string
		task        func() (any, error)
		expected    any
		expectedErr error
	}{
		{
			name: "task with struct type",
			task: func() (any, error) {
				return res, nil
			},
			expected:    res,
			expectedErr: nil,
		},
		{
			name: "task with pointer struct type",
			task: func() (any, error) {
				return &res, nil
			},
			expected:    &res,
			expectedErr: nil,
		},
		{
			name: "task with primitive type",
			task: func() (any, error) {
				return 42, nil
			},
			expected:    42,
			expectedErr: nil,
		},
		{
			name: "task with error",
			task: func() (any, error) {
				return 0, errors.New("error42")
			},
			expected:    0,
			expectedErr: errors.New("error42"),
		},
		{
			name: "catch panic",
			task: func() (any, error) {
				panic("error42")
			},
			expected:    nil,
			expectedErr: fmt.Errorf("panic: error42"),
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			task := Submit(c.task)
			wg := &sync.WaitGroup{}

			for i := 0; i < 5; i++ {
				wg.Add(1)

				go func() {
					defer wg.Done()

					actual, actualErr := task.Await()
					require.Equal(t, c.expectedErr, actualErr)
					assert.Equal(t, c.expected, actual)
				}()
			}

			wg.Wait()
		})
	}
}

func TestTask_Pipeline(t *testing.T) {
	cases := []struct {
		name        string
		task        func() ([]int, error)
		expected    []int
		expectedErr error
	}{
		{
			name: "simple pipeline",
			task: func() ([]int, error) {
				t1 := Submit(func() ([]int, error) {
					return []int{1}, nil
				})
				t2 := Submit(func() ([]int, error) {
					v, err := t1.Await()
					if err != nil {
						return nil, err
					}

					return append(v, 2), nil
				})
				t3 := Submit(func() ([]int, error) {
					v, err := t2.Await()
					if err != nil {
						return nil, err
					}

					return append(v, 3), nil
				})
				t4 := Submit(func() ([]int, error) {
					v, err := t3.Await()
					if err != nil {
						return nil, err
					}

					return append(v, 4), nil
				})
				t5 := Submit(func() ([]int, error) {
					v, err := t4.Await()
					if err != nil {
						return nil, err
					}

					return append(v, 5), nil
				})

				return t5.Await()
			},
			expected:    []int{1, 2, 3, 4, 5},
			expectedErr: nil,
		},
		{
			name: "pipeline with a shared task",
			task: func() ([]int, error) {
				t1 := Submit(func() ([]int, error) {
					return []int{1}, nil
				})
				t2 := Submit(func() ([]int, error) {
					v, err := t1.Await()
					if err != nil {
						return nil, err
					}

					return append(v, 2), nil
				})
				t3 := Submit(func() ([]int, error) {
					v, err := t1.Await()
					if err != nil {
						return nil, err
					}

					return append(v, 3), nil
				})

				v1, err := t2.Await()
				if err != nil {
					return nil, err
				}

				v2, err := t3.Await()
				if err != nil {
					return nil, err
				}

				return append(v1, v2...), nil
			},
			expected:    []int{1, 2, 1, 3},
			expectedErr: nil,
		},
		{
			name: "error in the second task",
			task: func() ([]int, error) {
				t1 := Submit(func() ([]int, error) {
					return []int{1}, nil
				})

				t2 := Submit(func() ([]int, error) {
					v, err := t1.Await()
					if err != nil {
						return nil, err
					}

					return v, errors.New("error42")
				})
				t3 := Submit(func() ([]int, error) {
					v, err := t2.Await()
					if err != nil {
						return nil, err
					}

					return append(v, 3), nil
				})

				return t3.Await()
			},
			expected:    nil,
			expectedErr: errors.New("error42"),
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			task := Submit(c.task)
			actual, actualErr := task.Await()
			require.Equal(t, c.expectedErr, actualErr)
			assert.Equal(t, c.expected, actual)
		})
	}
}

func TestActionTask_Await(t *testing.T) {
	cases := []struct {
		name     string
		task     func() error
		expected error
	}{
		{
			name: "no error example",
			task: func() error {
				return nil
			},
			expected: nil,
		},
		{
			name: "error example",
			task: func() error {
				return errors.New("error42")
			},
			expected: errors.New("error42"),
		},
	}

	for _, c := range cases {
		c := c

		t.Run(c.name, func(t *testing.T) {
			task := SubmitAction(c.task)
			wg := &sync.WaitGroup{}

			for i := 0; i < 5; i++ {
				wg.Add(1)

				go func() {
					defer wg.Done()

					actual := task.Await()
					assert.Equal(t, c.expected, actual)
				}()
			}

			wg.Wait()
		})
	}
}
