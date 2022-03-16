package api

import (
	"context"
	log "github.com/sirupsen/logrus"
)

type StatePrepairing struct {
}

func NewStatePrepairing() *StatePrepairing {
	return &StatePrepairing{}
}

func (s *StatePrepairing) Enter(_ context.Context, _ ...interface{}) error {
	log.Info("[preparing] enter")
	return nil
}

func (s *StatePrepairing) Exit(_ context.Context, _ ...interface{}) error {
	log.Info("[preparing] exit")
	return nil
}
