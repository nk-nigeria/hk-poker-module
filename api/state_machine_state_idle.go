package api

import (
	"context"
	log "github.com/sirupsen/logrus"
)

type StateIdle struct {
	StateBase
}

func NewIdleState(fn FireFn) *StateIdle {
	return &StateIdle{
		StateBase: StateBase{
			fireFn: fn,
		},
	}
}

func (s *StateIdle) Enter(ctx context.Context, _ ...interface{}) error {
	log.Info("[idle] enter")
	// setup idle timeout
	procPkg := GetProcessorPackagerFromContext(ctx)
	state := procPkg.GetState()
	state.SetUpCountDown(idleTimeout)

	return nil
}

func (s *StateIdle) Exit(_ context.Context, _ ...interface{}) error {
	log.Info("[idle] exit")
	return nil
}

func (s *StateIdle) Process(ctx context.Context, args ...interface{}) error {
	log.Infof("[idle] processing")
	procPkg := GetProcessorPackagerFromContext(ctx)
	state := procPkg.GetState()
	log.Info("state presences size ", state.GetPresenceSize())

	if state.GetPresenceSize() > 0 {
		s.Trigger(ctx, triggerMatching)
	}

	if state.GetRemainCountDown() < 0 {
		// Do finish here
		//s.Trigger(ctx, triggerFinish)
	}

	return nil
}
