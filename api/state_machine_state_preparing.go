package api

import (
	"context"
	log "github.com/ciaolink-game-platform/cgp-chinese-poker-module/pkg/log"
	pb "github.com/ciaolink-game-platform/cgp-chinese-poker-module/proto"
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
	procPkg := GetProcessorPackagerFromContext(ctx)
	state := procPkg.GetState()
	log.GetLogger().Info("state %v", state.Presences)
	state.SetUpCountDown(preparingTimeout)

	return nil
}

func (s *StatePreparing) Exit(_ context.Context, _ ...interface{}) error {
	log.GetLogger().Info("[preparing] exit")
	return nil
}

func (s *StatePreparing) Process(ctx context.Context, args ...interface{}) error {
	log.GetLogger().Info("[preparing] processing")
	procPkg := GetProcessorPackagerFromContext(ctx)
	state := procPkg.GetState()
	if remain := state.GetRemainCountDown(); remain > 0 {
		pbGameState := pb.UpdateGameState{
			State:     pb.GameState_GameStatePreparing,
			CountDown: int64(remain),
		}

		err := procPkg.GetProcessor().broadcastMessage(procPkg.GetLogger(), procPkg.GetDispatcher(), int64(pb.OpCodeUpdate_OPCODE_UPDATE_GAME_STATE), &pbGameState, nil, nil, true)
		if err != nil {
			log.GetLogger().Warn("broadcast message error %v", err)
		}
	} else {
		// check preparing condition
		log.GetLogger().Info("[preparing] preparing timeout check presence count")
		if state.IsReadyToPlay() {
			// change to play
			s.Trigger(ctx, triggerPreparingDone)
		} else {
			// change to wait
			s.Trigger(ctx, triggerPreparingFailed)
		}
	}

	return nil
}
