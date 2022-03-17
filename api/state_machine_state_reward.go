package api

import (
	"context"
	log "github.com/sirupsen/logrus"
)

type StateReward struct {
	StateBase
}

func NewStateReward(fn FireFn) *StateReward {
	return &StateReward{
		StateBase: StateBase{
			fireFn: fn,
		},
	}
}

func (s *StateReward) Enter(_ context.Context, _ ...interface{}) error {
	log.Info("[reward] enter")
	return nil
}

func (s *StateReward) Exit(_ context.Context, _ ...interface{}) error {
	log.Info("[reward] exit")
	return nil
}

func (s *StateReward) Process(ctx context.Context, args ...interface{}) error {
	return nil
}
