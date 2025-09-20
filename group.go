package async

import (
	"context"
	"sync"
)

type Group struct {
	wg     sync.WaitGroup
	cancel func()
}

func NewGroup(ctx context.Context) (*Group, context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	return &Group{cancel: cancel}, ctx
}

func (s *Group) Close() {
	s.cancel()
	s.wg.Wait()
}

func (s *Group) add() {
	s.wg.Add(1)
}

func (s *Group) done() {
	s.wg.Done()
}
