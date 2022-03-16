package api

import (
	"context"
	pb "github.com/ciaolink-game-platform/cgp-chinese-poker-module/proto"
	log "github.com/sirupsen/logrus"
	"time"
)

const preparingTimeout = time.Second * 10

type StatePreparing struct {
	fn FireFn
}

func NewStatePreparing(fn FireFn) *StatePreparing {
	return &StatePreparing{
		fn: fn,
	}
}

func (s *StatePreparing) Enter(_ context.Context, args ...interface{}) error {
	log.Info("[preparing] enter")
	state := GetState(args...)
	log.Infof("state %v", state.Presences)
	state.SetUpCountDown(preparingTimeout)

	return nil
}

func (s *StatePreparing) Exit(_ context.Context, _ ...interface{}) error {
	log.Info("[preparing] exit")
	return nil
}

func (s *StatePreparing) Process(ctx context.Context, args ...interface{}) error {
	log.Infof("[preparing] processing %v, %v", len(args), args[0])
	state := GetState(args...)
	if remain := state.GetRemainCountDown(); remain > 0 {
		pbGameState := pb.UpdateGameState{
			State:     pb.GameState_GameStatePreparing,
			CountDown: int64(remain),
		}

		err := GetProcessor(args...).broadcastMessage(GetLogger(args...), GetDispatcher(args...), int64(pb.OpCodeUpdate_OPCODE_UPDATE_GAME_STATE), &pbGameState, nil, nil, true)
		if err != nil {
			log.Warnf("broadcast message error %v", err)
		}
	} else {
		// check preparing condition
	}

	return nil
}
