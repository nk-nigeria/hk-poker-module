package api

import (
	"context"
	log "github.com/ciaolink-game-platform/cgp-chinese-poker-module/pkg/log"
	pb "github.com/ciaolink-game-platform/cgp-chinese-poker-module/proto"
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
	log.GetLogger().Info("[matching] enter")
	procPkg := GetProcessorPackagerFromContext(ctx)
	state := procPkg.GetState()

	procPkg.GetProcessor().notifyUpdateGameState(
		state,
		procPkg.GetLogger(),
		procPkg.GetDispatcher(),
		&pb.UpdateGameState{
			State: pb.GameState_GameStateMatching,
		},
	)

	return nil
}

func (s *StateMatching) Exit(_ context.Context, _ ...interface{}) error {
	log.GetLogger().Info("[matching] exit")
	return nil
}

func (s *StateMatching) Process(ctx context.Context, args ...interface{}) error {
	log.GetLogger().Info("[matching] processing")
	procPkg := GetProcessorPackagerFromContext(ctx)
	state := procPkg.GetState()
	presenceCount := state.GetPresenceSize()
	if state.IsReadyToPlay() {
		s.Trigger(ctx, triggerPresenceReady)
	} else if presenceCount <= 0 {
		s.Trigger(ctx, triggerIdle)
	} else {
		log.GetLogger().Info("state idle presences size ", presenceCount)
	}

	return nil
}
