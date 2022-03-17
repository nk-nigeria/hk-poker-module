package api

import (
	"context"
	pb "github.com/ciaolink-game-platform/cgp-chinese-poker-module/proto"
	log "github.com/sirupsen/logrus"
)

type StatePlay struct {
	StateBase
}

func NewStatePlay(fn FireFn) *StatePlay {
	return &StatePlay{
		StateBase: StateBase{
			fireFn: fn,
		},
	}
}

func (s *StatePlay) Enter(ctx context.Context, agrs ...interface{}) error {
	log.Info("[play] enter")
	procPkg := GetProcessorPackagerFromContext(ctx)
	state := procPkg.GetState()
	// Setup count down
	state.SetUpCountDown(playTimeout)
	// New game here
	procPkg.GetProcessor().processNewGame(procPkg.GetLogger(), procPkg.GetDispatcher(), state)

	return nil
}

func (s *StatePlay) Exit(_ context.Context, _ ...interface{}) error {
	log.Info("[play] exit")
	return nil
}

func (s *StatePlay) Process(ctx context.Context, args ...interface{}) error {
	log.Infof("[play] processing")
	procPkg := GetProcessorPackagerFromContext(ctx)
	state := procPkg.GetState()
	if remain := state.GetRemainCountDown(); remain > 0 {
		log.Infof("[play] not timeout %v", remain)
		messages := procPkg.GetMessages()
		processor := procPkg.GetProcessor()
		logger := procPkg.GetLogger()
		dispatcher := procPkg.GetDispatcher()
		for _, message := range messages {
			switch pb.OpCodeRequest(message.GetOpCode()) {
			case pb.OpCodeRequest_OPCODE_REQUEST_COMBINE_CARDS:
				processor.combineCard(logger, dispatcher, state, message)
			case pb.OpCodeRequest_OPCODE_REQUEST_SHOW_CARDS:
				processor.showCard(logger, dispatcher, state, message)
			}
		}
	} else {
		log.Infof("[play] timeout reach %v", remain)
		s.Trigger(ctx, triggerPlayTimeout, state)
	}
	return nil
}
