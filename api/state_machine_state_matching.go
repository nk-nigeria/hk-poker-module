package api

import (
	"context"
	log "github.com/sirupsen/logrus"
)

type StateMatching struct {
	StateBase
}

func NewStateMatching(fn FireFn) *StateMatching {
	return &StateMatching{
		StateBase: StateBase{
			fireFn: fn,
		},
	}
}

func (s *StateMatching) Enter(ctx context.Context, _ ...interface{}) error {
	log.Info("[matching] enter")

	return nil
}

func (s *StateMatching) Exit(_ context.Context, _ ...interface{}) error {
	log.Info("[matching] exit")
	return nil
}

func (s *StateMatching) Process(ctx context.Context, args ...interface{}) error {
	log.Infof("[matching] processing")
	procPkg := GetProcessorPackagerFromContext(ctx)
	state := procPkg.GetState()
	presenceCount := state.GetPresenceSize()
	if state.IsReadyToPlay() {
		s.Trigger(ctx, triggerPresenceReady)
	} else if presenceCount <= 0 {
		s.Trigger(ctx, triggerIdle)
	} else {
		log.Info("state idle presences size ", presenceCount)
	}

	return nil
}
