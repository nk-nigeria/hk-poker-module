package api

import (
	"context"
	_ "github.com/filecoin-project/go-statemachine"
	_ "github.com/ipfs/go-datastore"
	"github.com/qmuntal/stateless"
	"time"
)

const (
	stateWait      = "Wait"
	statePreparing = "Preparing"
	statePlay      = "Play"
	stateReward    = "Reward"
	stateFinish    = "Finish"
)

const (
	triggerPresenceReady   = "GamePresenceReady"
	triggerPreparingDone   = "GamePrepairingDone"
	triggerPreparingFailed = "GamePrepairingFailed"
	triggerPlayTimeout     = "GamePlayTimeout"
	triggerPlayCombineAll  = "GamePlayCombineAll"
	triggerRewardTimeout   = "GameRewardTimeout"
	triggerNoOne           = "GameNoOne"

	triggerProcessWait      = "GameProcessWait"
	triggerProcessPreparing = "GameProcessPrepairing"
	triggerProcessPlay      = "GameProcessPlay"
	triggerProcessReward    = "GameProcessReward"
)

const (
	preparingTimeout = time.Second * 10
	playTimeout      = time.Second * 60
)

type GameStateMachine struct {
	state *stateless.StateMachine
}

func (m *GameStateMachine) configure() {
	fireCtx := m.state.FireCtx
	wait := NewStateWait(fireCtx)
	m.state.Configure(stateWait).
		OnEntry(wait.Enter).
		OnExit(wait.Exit).
		InternalTransition(triggerProcessWait, wait.Process).
		Permit(triggerPresenceReady, statePreparing).
		Permit(triggerNoOne, stateFinish)

	preparing := NewStatePreparing(fireCtx)
	m.state.Configure(statePreparing).
		OnEntry(preparing.Enter).
		OnExit(preparing.Exit).
		InternalTransition(triggerProcessPreparing, preparing.Process).
		Permit(triggerPreparingDone, statePlay).
		Permit(triggerPreparingFailed, stateWait)

	play := NewStatePlay(fireCtx)
	m.state.Configure(statePlay).
		OnEntry(play.Enter).
		OnExit(play.Exit).
		InternalTransition(triggerProcessPlay, play.Process).
		Permit(triggerPlayTimeout, stateReward).
		Permit(triggerPlayCombineAll, stateReward)

	reward := NewStateReward(fireCtx)
	m.state.Configure(stateReward).
		OnEntry(reward.Enter).
		OnExit(reward.Exit).
		InternalTransition(triggerProcessReward, reward.Process).
		Permit(triggerRewardTimeout, statePreparing)

	m.state.ToGraph()
}

func (m *GameStateMachine) FireProcessEvent(ctx context.Context, args ...interface{}) error {
	var trigger stateless.State
	switch m.state.MustState() {
	case stateWait:
		trigger = triggerProcessWait
	case statePreparing:
		trigger = triggerProcessPreparing
	case statePlay:
		trigger = triggerProcessPlay
	case stateReward:
		trigger = triggerProcessReward
	}
	return m.state.FireCtx(ctx, trigger, args...)
}

func (m *GameStateMachine) MustState() stateless.State {
	return m.state.MustState()
}

func NewGameStateMachine() *GameStateMachine {
	gs := &GameStateMachine{
		state: stateless.NewStateMachine(stateWait),
	}

	gs.configure()

	return gs
}

type FireFn func(ctx context.Context, trigger stateless.Trigger, args ...interface{}) error
type StateBase struct {
	fireFn FireFn
}

func (s *StateBase) Trigger(ctx context.Context, trigger stateless.Trigger, args ...interface{}) error {
	return s.fireFn(ctx, trigger, args...)
}
