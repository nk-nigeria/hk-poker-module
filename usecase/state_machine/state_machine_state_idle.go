package state_machine

import (
	"context"

	pb "github.com/nk-nigeria/cgp-common/proto"
	"github.com/nk-nigeria/hk-poker-module/api/presenter"
	"github.com/nk-nigeria/hk-poker-module/entity"
	log "github.com/nk-nigeria/hk-poker-module/pkg/log"
	"github.com/nk-nigeria/hk-poker-module/pkg/packager"
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
	state.SetUpCountDown(entity.GameStateDuration[state.GetGameState()])

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
			State: pb.GameState_GAME_STATE_IDLE,
		},
	)

	return nil
}

func (s *StateIdle) Exit(_ context.Context, _ ...interface{}) error {
	log.GetLogger().Info("[idle] exit")
	return nil
}

func (s *StateIdle) Process(ctx context.Context, args ...interface{}) error {
	// log.GetLogger().Info("[idle] processing")
	procPkg := packager.GetProcessorPackagerFromContext(ctx)
	state := procPkg.GetState()
	// log.GetLogger().Info("state presences size %v", state.GetPresenceSize())

	if state.GetPrecenseNotBotCount() > 0 {
		s.Trigger(ctx, triggerMatching)
		return nil
	}

	if remain := state.GetRemainCountDown(); remain < 0 {
		log.GetLogger().Info("[idle] idle timeout => exit")
		s.Trigger(ctx, triggerNoOne)
		return presenter.ErrGameFinish
	}

	return nil
}
