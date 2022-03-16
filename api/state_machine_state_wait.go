package api

import (
	"context"
	log "github.com/sirupsen/logrus"
)

type StateWait struct {
	fn FireFn
}

func NewStateWait(fn FireFn) *StateWait {
	return &StateWait{
		fn: fn,
	}
}

func (s *StateWait) Enter(_ context.Context, _ ...interface{}) error {
	log.Info("[wait] enter")
	return nil
}

func (s *StateWait) Exit(_ context.Context, _ ...interface{}) error {
	log.Info("[wait] exit")
	return nil
}

func (s *StateWait) Process(ctx context.Context, args ...interface{}) error {
	log.Infof("[wait] processing %v, %v", len(args), args[0])
	state := GetState(args...)
	log.Info("state presences size ", state.Presences.Size())
	if state.Presences.Size() >= state.MinPresences {
		s.fn(triggerPresenceReady, state)
	}

	return nil
}
