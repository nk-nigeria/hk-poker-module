package state_machine

import (
	"context"
	"time"

	"github.com/nakamaFramework/cgp-chinese-poker-module/entity"
	log "github.com/nakamaFramework/cgp-chinese-poker-module/pkg/log"
	"github.com/nakamaFramework/cgp-chinese-poker-module/pkg/packager"
	pb "github.com/nakamaFramework/cgp-common/proto"
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
	log.GetLogger().Info("[play] enter")
	procPkg := packager.GetProcessorPackagerFromContext(ctx)
	state := procPkg.GetState()
	state.MatchCount++
	// Setup count down
	state.SetUpCountDown(entity.GameStateDuration[state.GetGameState()])
	state.SetupMatchPresence()
	procPkg.GetProcessor().ProcessNewGame(procPkg.GetContext(),
		procPkg.GetNK(),
		procPkg.GetDb(),
		procPkg.GetLogger(),
		procPkg.GetDispatcher(), state)

	time.Sleep(time.Millisecond * 200)
	state.DelayForDeclare.Setup(1*time.Second, entity.TickRate)
	procPkg.GetProcessor().NotifyUpdateGameState(
		state,
		procPkg.GetLogger(),
		procPkg.GetDispatcher(),
		&pb.UpdateGameState{
			State:     pb.GameState_GameStatePlay,
			CountDown: int64(state.GetRemainCountDown()),
		},
	)
	return nil
}

func (s *StatePlay) Exit(_ context.Context, _ ...interface{}) error {
	log.GetLogger().Info("[play] exit")
	return nil
}

func (s *StatePlay) Process(ctx context.Context, args ...interface{}) error {
	// log.GetLogger().Info("[play] processing")
	procPkg := packager.GetProcessorPackagerFromContext(ctx)
	state := procPkg.GetState()
	remain := state.GetRemainCountDown()
	if remain > 0 || !state.DelayForDeclare.Timeout() {
		messages := procPkg.GetMessages()
		processor := procPkg.GetProcessor()
		logger := procPkg.GetLogger()
		dispatcher := procPkg.GetDispatcher()
		processor.ProcessGame(procPkg.GetContext(), logger, procPkg.GetNK(), procPkg.GetDb(), dispatcher, messages, state)
		if remain <= 0 {
			state.DelayForDeclare.Loop()
		}
		// Check all user show card
		if state.GetShowCardCount() >= state.GetPlayingCount() {
			s.Trigger(ctx, triggerPlayCombineAll)
		}
	} else {
		log.GetLogger().Info("[play] timeout reach %v", remain)
		s.Trigger(ctx, triggerPlayTimeout)
	}
	return nil
}
