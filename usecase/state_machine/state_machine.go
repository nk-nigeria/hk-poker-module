package state_machine

import (
	"context"

	"github.com/nakamaFramework/cgp-chinese-poker-module/entity"
	"github.com/nakamaFramework/cgp-chinese-poker-module/pkg/packager"
	pb "github.com/nakamaFramework/cgp-common/proto"
	"github.com/qmuntal/stateless"
)

const (
	StateInitType      = pb.GameState_GameStateUnknown // Only for initialize
	StateIdleType      = pb.GameState_GameStateIdle
	StateMatchingType  = pb.GameState_GameStateMatching
	StatePreparingType = pb.GameState_GameStatePreparing
	StatePlayType      = pb.GameState_GameStatePlay
	StateRewardType    = pb.GameState_GameStateReward
	StateFinishType    = pb.GameState_GameStateFinish
)

const (
	triggerIdle            = "GameIdle"
	triggerMatching        = "GameMatching"
	triggerPresenceReady   = "GamePresenceReady"
	triggerPreparingDone   = "GamePreparingDone"
	triggerPreparingFailed = "GamePreparingFailed"
	triggerPlayTimeout     = "GamePlayTimeout"
	triggerPlayCombineAll  = "GamePlayCombineAll"
	triggerRewardTimeout   = "GameRewardTimeout"
	triggerNoOne           = "GameNoOne"

	triggerProcess = "GameProcess"
)

type Machine struct {
	state *stateless.StateMachine
}

func (m *Machine) configure() {
	fireCtx := m.state.FireCtx

	// init state
	m.state.Configure(StateInitType).
		Permit(triggerIdle, StateIdleType)
	m.state.OnTransitioning(func(ctx context.Context, t stateless.Transition) {
		procPkg := packager.GetProcessorPackagerFromContext(ctx)
		state := procPkg.GetState()
		var ok bool
		state.Label.GameState, ok = t.Destination.(pb.GameState)
		if !ok {
			return
		}
		if procPkg.GetDispatcher() != nil {
			labelJson, _ := entity.DefaultMarshaler.Marshal(state.Label)
			procPkg.GetDispatcher().MatchLabelUpdate(string(labelJson))
		}
	})

	// idle state: wait for first user, check no one and timeout
	idle := NewIdleState(fireCtx)
	m.state.Configure(StateIdleType).
		OnEntry(idle.Enter).
		OnExit(idle.Exit).
		InternalTransition(triggerProcess, idle.Process).
		Permit(triggerMatching, StateMatchingType).
		Permit(triggerNoOne, StateFinishType)

	// matching state: wait for reach min user => switch to preparing, check no one and timeout => switch to idle
	matching := NewStateMatching(fireCtx)
	m.state.Configure(StateMatchingType).
		OnEntry(matching.Enter).
		OnExit(matching.Exit).
		InternalTransition(triggerProcess, matching.Process).
		Permit(triggerPresenceReady, StatePreparingType).
		Permit(triggerIdle, StateIdleType)

	// preparing state: wait for reach min user in duration => switch to play, check not enough and timeout => switch to idle
	preparing := NewStatePreparing(fireCtx)
	m.state.Configure(StatePreparingType).
		OnEntry(preparing.Enter).
		OnExit(preparing.Exit).
		InternalTransition(triggerProcess, preparing.Process).
		Permit(triggerPreparingDone, StatePlayType).
		Permit(triggerPreparingFailed, StateMatchingType)

	// playing state: wait for all user show card or timeout =>
	//  switch to reward
	play := NewStatePlay(fireCtx)
	m.state.Configure(StatePlayType).
		OnEntry(play.Enter).
		OnExit(play.Exit).
		InternalTransition(triggerProcess, play.Process).
		Permit(triggerPlayTimeout, StateRewardType).
		Permit(triggerPlayCombineAll, StateRewardType)

	// reward state: wait for reward timeout => switch to
	reward := NewStateReward(fireCtx)
	m.state.Configure(StateRewardType).
		OnEntry(reward.Enter).
		OnExit(reward.Exit).
		InternalTransition(triggerProcess, reward.Process).
		Permit(triggerRewardTimeout, StateMatchingType)

	m.state.ToGraph()
}

func (m *Machine) FireProcessEvent(ctx context.Context, args ...interface{}) error {
	return m.state.FireCtx(ctx, triggerProcess, args...)
}

func (m *Machine) MustState() stateless.State {
	return m.state.MustState()
}

func (m *Machine) GetPbState() pb.GameState {
	switch m.state.MustState() {
	case StateIdleType:
		return pb.GameState_GameStateIdle
	case StateMatchingType:
		return pb.GameState_GameStateMatching
	case StatePreparingType:
		return pb.GameState_GameStatePreparing
	case StatePlayType:
		return pb.GameState_GameStatePlay
	case StateRewardType:
		return pb.GameState_GameStateReward
	default:
		return pb.GameState_GameStateUnknown
	}
}

func NewGameStateMachine() UseCase {
	gs := &Machine{
		state: stateless.NewStateMachine(StateInitType),
	}

	gs.configure()

	return gs
}

func (m *Machine) IsPlayingState() bool {
	return m.MustState() == StatePlayType
}

func (m *Machine) IsReward() bool {
	return m.MustState() == StateRewardType
}

func (m *Machine) Trigger(ctx context.Context, trigger stateless.Trigger, args ...interface{}) error {
	return m.state.FireCtx(ctx, trigger, args...)
}

func (m *Machine) TriggerIdle(ctx context.Context, args ...interface{}) error {
	return m.state.FireCtx(ctx, triggerIdle, args...)
}

type FireFn func(ctx context.Context, trigger stateless.Trigger, args ...interface{}) error
type StateBase struct {
	fireFn FireFn
}

func (s *StateBase) Trigger(ctx context.Context, trigger stateless.Trigger, args ...interface{}) error {
	return s.fireFn(ctx, trigger, args...)
}
