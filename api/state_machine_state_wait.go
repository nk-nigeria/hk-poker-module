package api

import (
	"context"
	log "github.com/sirupsen/logrus"
)

type StateWait struct {
}

func NewStateWait() *StateWait {
	return &StateWait{}
}

func (s *StateWait) Enter(_ context.Context, _ ...interface{}) error {
	log.Info("[wait] enter")
	return nil
}

func (s *StateWait) Exit(_ context.Context, _ ...interface{}) error {
	log.Info("[wait] exit")
	return nil
}

func (s *StateWait) Process(ctx context.Context, i ...interface{}) error {
	return nil
}
