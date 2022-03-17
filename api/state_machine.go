package api

import (
	"context"
	pb "github.com/ciaolink-game-platform/cgp-chinese-poker-module/proto"
	"github.com/qmuntal/stateless"
	"time"
)

const (
	stateInit      = "Init" // Only for initialize
	stateIdle      = "Idle"
	stateMatching  = "Matching"
	statePreparing = "Preparing"
	statePlay      = "Play"
	stateReward    = "Reward"
	stateFinish    = "Finish"
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

	triggerProcessIdle      = "GameProcessIdle"
	triggerProcessMatching  = "GameProcessMatching"
	triggerProcessPreparing = "GameProcessPreparing"
	triggerProcessPlay      = "GameProcessPlay"
	triggerProcessReward    = "GameProcessReward"
)

const (
	idleTimeout      = time.Second * 15
	preparingTimeout = time.Second * 10
	playTimeout      = time.Second * 60
	rewardTimeout    = time.Second * 30
)

type GameStateMachine struct {
	state *stateless.StateMachine
}

func (m *GameStateMachine) configure() {
	fireCtx := m.state.FireCtx

	// init state
	m.state.Configure(stateInit).
		Permit(triggerIdle, stateIdle)

	// idle state: wait for first user, check no one and timeout
	idle := NewIdleState(fireCtx)
	m.state.Configure(stateIdle).
		OnEntry(idle.Enter).
		OnExit(idle.Exit).
		InternalTransition(triggerProcessIdle, idle.Process).
		Permit(triggerMatching, stateMatching).
		Permit(triggerNoOne, stateFinish)

	// matching state: wait for reach min user => switch to preparing, check no one and timeout => switch to idle
	matching := NewStateMatching(fireCtx)
	m.state.Configure(stateMatching).
		OnEntry(matching.Enter).
		OnExit(matching.Exit).
		InternalTransition(triggerProcessMatching, matching.Process).
		Permit(triggerPresenceReady, statePreparing).
		Permit(triggerIdle, stateIdle)

	// preparing state: wait for reach min user in duration => switch to play, check not enough and timeout => switch to idle
	preparing := NewStatePreparing(fireCtx)
	m.state.Configure(statePreparing).
		OnEntry(preparing.Enter).
		OnExit(preparing.Exit).
		InternalTransition(triggerProcessPreparing, preparing.Process).
		Permit(triggerPreparingDone, statePlay).
		Permit(triggerPreparingFailed, stateMatching)

	// playing state: wait for all user show card or timeout => switch to reward
	play := NewStatePlay(fireCtx)
	m.state.Configure(statePlay).
		OnEntry(play.Enter).
		OnExit(play.Exit).
		InternalTransition(triggerProcessPlay, play.Process).
		Permit(triggerPlayTimeout, stateReward).
		Permit(triggerPlayCombineAll, stateReward)

	// reward state: wait for reward timeout => switch to
	reward := NewStateReward(fireCtx)
	m.state.Configure(stateReward).
		OnEntry(reward.Enter).
		OnExit(reward.Exit).
		InternalTransition(triggerProcessReward, reward.Process).
		Permit(triggerRewardTimeout, stateMatching)

	m.state.ToGraph()
}

func (m *GameStateMachine) FireProcessEvent(ctx context.Context, args ...interface{}) error {
	var trigger stateless.State
	switch m.state.MustState() {
	case stateIdle:
		trigger = triggerProcessIdle
	case stateMatching:
		trigger = triggerProcessMatching
	case statePreparing:
		trigger = triggerProcessPreparing
	case statePlay:
		trigger = triggerProcessPlay
	case stateReward:
		trigger = triggerProcessReward
	default:
		return nil
	}
	return m.state.FireCtx(ctx, trigger, args...)
}

func (m *GameStateMachine) MustState() stateless.State {
	return m.state.MustState()
}

func (m *GameStateMachine) GetPbState() pb.GameState {
	switch m.state.MustState() {
	case stateIdle:
		return pb.GameState_GameStateIdle
	case stateMatching:
		return pb.GameState_GameStateMatching
	case statePreparing:
		return pb.GameState_GameStatePreparing
	case statePlay:
		return pb.GameState_GameStatePlay
	case stateReward:
		return pb.GameState_GameStateReward
	default:
		return pb.GameState_GameStateUnknown
	}
}

func NewGameStateMachine() *GameStateMachine {
	gs := &GameStateMachine{
		state: stateless.NewStateMachine(stateInit),
	}

	gs.configure()

	return gs
}

func (m *GameStateMachine) Trigger(ctx context.Context, trigger stateless.Trigger, args ...interface{}) error {
	return m.state.FireCtx(ctx, trigger, args...)
}

type FireFn func(ctx context.Context, trigger stateless.Trigger, args ...interface{}) error
type StateBase struct {
	fireFn FireFn
}

func (s *StateBase) Trigger(ctx context.Context, trigger stateless.Trigger, args ...interface{}) error {
	return s.fireFn(ctx, trigger, args...)
}
