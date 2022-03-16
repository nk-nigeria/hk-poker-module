package api

import (
	_ "github.com/filecoin-project/go-statemachine"
	_ "github.com/ipfs/go-datastore"
	"github.com/qmuntal/stateless"
)

const (
	stateWait       = "Wait"
	statePrepairing = "Preparing"
	statePlay       = "Play"
	stateReward     = "Reward"
	stateFinish     = "Finish"
)

const (
	triggerPresenceReady    = "GamePresenceReady"
	triggerPrepairingDone   = "GamePrepairingDone"
	triggerPrepairingFailed = "GamePrepairingFailed"
	triggerPlayTimeout      = "GamePlayTimeout"
	triggerPlayCombineAll   = "GamePlayCombineAll"
	triggerRewardTimeout    = "GameRewardTimeout"
	triggerNoOne            = "GameNoOne"

	triggerProcessWait      = "GameProcessWait"
	triggerProcessPreparing = "GameProcessPrepairing"
	triggerProcessPlay      = "GameProcessPlay"
	triggerProcessReward    = "GameProcessReward"
)

type GameStateMachine struct {
	state *stateless.StateMachine
}

func (m *GameStateMachine) configure() {
	wait := NewStateWait()
	m.state.Configure(stateWait).
		OnEntry(wait.Enter).
		OnExit(wait.Exit).
		InternalTransition(triggerProcessWait, wait.Process).
		Permit(triggerPresenceReady, statePrepairing).
		Permit(triggerNoOne, stateFinish)

	preparing := NewStatePrepairing()
	m.state.Configure(statePrepairing).
		OnEntry(preparing.Enter).
		OnExit(preparing.Exit).
		InternalTransition(triggerProcessPreparing, preparing.Process).
		Permit(triggerPrepairingDone, statePlay).
		Permit(triggerPrepairingFailed, stateWait)

	run := NewStateRun()
	m.state.Configure(statePlay).
		OnEntry(run.Enter).
		OnExit(run.Exit).
		InternalTransition(triggerProcessPlay, run.Process).
		Permit(triggerPlayTimeout, stateReward).
		Permit(triggerPlayCombineAll, stateReward)

	reward := NewStateReward()
	m.state.Configure(stateReward).
		OnEntry(reward.Enter).
		OnExit(reward.Exit).
		InternalTransition(triggerProcessReward, reward.Process).
		Permit(triggerRewardTimeout, statePrepairing)

	m.state.ToGraph()
}

func (m *GameStateMachine) Fire(trigger stateless.Trigger, args ...interface{}) error {
	return m.state.Fire(trigger, args)
}

func (m *GameStateMachine) FireProcessEvent(args ...interface{}) error {
	var trigger stateless.State
	switch m.state.MustState() {
	case stateWait:
		trigger = triggerProcessWait
	case statePrepairing:
		trigger = triggerProcessPreparing
	case statePlay:
		trigger = triggerProcessPlay
	case stateReward:
		trigger = triggerProcessReward
	}
	return m.state.Fire(trigger, args)
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
