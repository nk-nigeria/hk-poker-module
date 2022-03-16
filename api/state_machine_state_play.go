package api

import (
	"context"
	log "github.com/sirupsen/logrus"
)

type StatePlay struct {
}

func NewStatePlay() *StatePlay {
	return &StatePlay{}
}

func (s *StatePlay) Enter(_ context.Context, _ ...interface{}) error {
	log.Info("[play] enter")
	return nil
}

func (s *StatePlay) Exit(_ context.Context, _ ...interface{}) error {
	log.Info("[play] exit")
	return nil
}

func (s *StatePlay) Process(ctx context.Context, args ...interface{}) error {
	return nil
}
