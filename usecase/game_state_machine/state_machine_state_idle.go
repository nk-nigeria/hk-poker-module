package game_state_machine

import (
	"context"
	"github.com/ciaolink-game-platform/cgp-chinese-poker-module/api/presenter"
	log "github.com/ciaolink-game-platform/cgp-chinese-poker-module/pkg/log"
	"github.com/ciaolink-game-platform/cgp-chinese-poker-module/pkg/packager"
	pb "github.com/ciaolink-game-platform/cgp-chinese-poker-module/proto"
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
	log.GetLogger().Info("[idle] enter")
	// setup idle timeout
	procPkg := packager.GetProcessorPackagerFromContext(ctx)
	state := procPkg.GetState()
	state.SetUpCountDown(idleTimeout)

	dispatcher := procPkg.GetDispatcher()
	if dispatcher == nil {
		log.GetLogger().Warn("missing dispatcher don't broadcast")
		return nil
	}

	procPkg.GetProcessor().NotifyUpdateGameState(
		state,
		procPkg.GetLogger(),
		procPkg.GetDispatcher(),
		&pb.UpdateGameState{
			State: pb.GameState_GameStateIdle,
		},
	)

	return nil
}

func (s *StateIdle) Exit(_ context.Context, _ ...interface{}) error {
	log.GetLogger().Info("[idle] exit")
	return nil
}

func (s *StateIdle) Process(ctx context.Context, args ...interface{}) error {
	log.GetLogger().Info("[idle] processing")
	procPkg := packager.GetProcessorPackagerFromContext(ctx)
	state := procPkg.GetState()
	log.GetLogger().Info("state presences size %v", state.GetPresenceSize())

	if state.GetPresenceSize() > 0 {
		s.Trigger(ctx, triggerMatching)
	}

	if remain := state.GetRemainCountDown(); remain < 0 {
		// Do finish here
		//s.Trigger(ctx, triggerFinish)
		log.GetLogger().Info("[idle] idle timeout => exit")
		return presenter.ErrGameFinish
	} else {
		log.GetLogger().Info("[idle] idle timeout remain %v", remain)
	}

	return nil
}
