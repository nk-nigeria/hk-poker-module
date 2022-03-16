package api

import (
	"context"
	log "github.com/sirupsen/logrus"
)

type StateRun struct {
}

func NewStateRun() *StateRun {
	return &StateRun{}
}

func (s *StateRun) Enter(_ context.Context, _ ...interface{}) error {
	log.Info("[run] enter")
	return nil
}

func (s *StateRun) Exit(_ context.Context, _ ...interface{}) error {
	log.Info("[run] exit")
	return nil
}

func (s *StateRun) Process(ctx context.Context, args ...interface{}) error {
	return nil
}
