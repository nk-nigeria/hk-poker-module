package api

import (
	"context"
	log "github.com/sirupsen/logrus"
)

type StateWait struct {
	StateBase
}

func NewStateWait(fn FireFn) *StateWait {
	return &StateWait{
		StateBase: StateBase{
			fireFn: fn,
		},
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
	log.Infof("[wait] processing")
	procPkg := GetProcessorPackagerFromContext(ctx)
	state := procPkg.GetState()
	log.Info("state presences size ", state.GetPresenceSize())
	if state.GetPresenceSize() >= state.MinPresences {
		s.Trigger(ctx, triggerPresenceReady)
	} else {
		//TODO: check finish timeout

	}

	return nil
}
