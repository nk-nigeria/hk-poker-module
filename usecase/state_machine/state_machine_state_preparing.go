package state_machine

import (
	"context"
	"strings"

	pb "github.com/nk-nigeria/cgp-common/proto"
	"github.com/nk-nigeria/hk-poker-module/entity"
	log "github.com/nk-nigeria/hk-poker-module/pkg/log"
	"github.com/nk-nigeria/hk-poker-module/pkg/packager"
)

type StatePreparing struct {
	StateBase
}

func NewStatePreparing(fn FireFn) *StatePreparing {
	return &StatePreparing{
		StateBase: StateBase{
			fireFn: fn,
		},
	}
}

func (s *StatePreparing) Enter(ctx context.Context, args ...interface{}) error {
	log.GetLogger().Info("[preparing] enter")
	procPkg := packager.GetProcessorPackagerFromContext(ctx)
	state := procPkg.GetState()
	log.GetLogger().Info("state %v", state.Presences)
	state.SetUpCountDown(entity.GameStateDuration[state.GetGameState()])
	// remove all user not interact 2 game continue
	listPrecense := state.GetPresenceNotInteract(1)
	if len(listPrecense) > 0 {
		listUserId := make([]string, len(listPrecense))
		for _, p := range listPrecense {
			listUserId = append(listUserId, p.GetUserId())
		}
		procPkg.GetLogger().Info("Kick %d user from math %s",
			len(listPrecense), strings.Join(listUserId, ","))
		state.AddLeavePresence(listPrecense...)
	}
	procPkg.GetProcessor().ProcessApplyPresencesLeave(ctx,
		procPkg.GetLogger(),
		procPkg.GetNK(),
		procPkg.GetDb(),
		procPkg.GetDispatcher(),
		state,
	)

	procPkg.GetProcessor().NotifyUpdateGameState(
		state,
		procPkg.GetLogger(),
		procPkg.GetDispatcher(),
		&pb.UpdateGameState{
			State:     pb.GameState_GAME_STATE_PREPARING,
			CountDown: int64(state.GetRemainCountDown()),
		},
	)

	return nil
}

func (s *StatePreparing) Exit(_ context.Context, _ ...interface{}) error {
	log.GetLogger().Info("[preparing] exit")
	return nil
}

func (s *StatePreparing) Process(ctx context.Context, args ...interface{}) error {
	// log.GetLogger().Info("[preparing] processing")
	procPkg := packager.GetProcessorPackagerFromContext(ctx)
	state := procPkg.GetState()
	// if state.GetPrecenseNotBotCount() == 0 {
	// 	s.Trigger(ctx, triggerPreparingFailed)
	// 	return nil
	// }
	msgs := procPkg.GetMessages()
	if len(msgs) > 0 {
		procPkg.GetProcessor().ProcessMessageUser(procPkg.GetContext(),
			procPkg.GetLogger(),
			procPkg.GetNK(),
			procPkg.GetDb(),
			procPkg.GetDispatcher(),
			msgs, state)
	}
	if remain := state.GetRemainCountDown(); remain >= 1 {
		if state.IsNeedNotifyCountDown() {
			procPkg.GetProcessor().NotifyUpdateGameState(
				state,
				procPkg.GetLogger(),
				procPkg.GetDispatcher(),
				&pb.UpdateGameState{
					State:     pb.GameState_GAME_STATE_PREPARING,
					CountDown: int64(remain),
				},
			)

			state.SetLastCountDown(remain)
		}
		return nil

	}
	if !state.IsReadyToPlay() {
		s.Trigger(ctx, triggerPreparingFailed)
		return nil
	}
	s.Trigger(ctx, triggerPreparingDone)
	return nil
}
