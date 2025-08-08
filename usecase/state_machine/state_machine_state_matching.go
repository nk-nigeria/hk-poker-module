package state_machine

import (
	"context"

	pb "github.com/nk-nigeria/cgp-common/proto"
	"github.com/nk-nigeria/hk-poker-module/api/presenter"
	"github.com/nk-nigeria/hk-poker-module/entity"
	log "github.com/nk-nigeria/hk-poker-module/pkg/log"
	"github.com/nk-nigeria/hk-poker-module/pkg/packager"
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
	procPkg := packager.GetProcessorPackagerFromContext(ctx)
	state := procPkg.GetState()
	state.SetUpCountDown(entity.GameStateDuration[state.GetGameState()])
	procPkg.GetLogger().Info("apply leave presence")

	procPkg.GetProcessor().ProcessApplyPresencesLeave(
		procPkg.GetContext(),
		procPkg.GetLogger(),
		procPkg.GetNK(),
		procPkg.GetDb(),
		procPkg.GetDispatcher(),
		state)
	procPkg.GetProcessor().NotifyUpdateGameState(
		state,
		procPkg.GetLogger(),
		procPkg.GetDispatcher(),
		&pb.UpdateGameState{
			State: pb.GameState_GAME_STATE_MATCHING,
		},
	)
	return nil
}

func (s *StateMatching) Exit(_ context.Context, _ ...interface{}) error {
	log.GetLogger().Info("[matching] exit")
	return nil
}

func (s *StateMatching) Process(ctx context.Context, args ...interface{}) error {
	// log.GetLogger().Info("[matching] processing")
	procPkg := packager.GetProcessorPackagerFromContext(ctx)
	state := procPkg.GetState()
	remain := state.GetRemainCountDown()
	if state.GetPrecenseBotCount() == state.GetPresenceSize() {
		// s.Trigger(ctx, triggerIdle)
		// return nil
		return presenter.ErrGameFinish
	}
	if remain > 0 {
		return nil
	}

	log.GetLogger().WithField("count presence", state.GetPresenceSize()).Info("[matching] processing")
	presenceCount := state.GetPresenceSize()
	if presenceCount == state.GetPrecenseBotCount() {
		s.Trigger(ctx, triggerIdle)
		return nil
	}
	if presenceCount >= state.MinPresences {
		s.Trigger(ctx, triggerPresenceReady)
	} else if presenceCount == state.GetPrecenseBotCount() {
		s.Trigger(ctx, triggerNoOne)
	} else {
		s.Trigger(ctx, triggerIdle)
	}

	return nil
}
