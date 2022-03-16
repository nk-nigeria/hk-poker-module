package api

import (
	_ "github.com/filecoin-project/go-statemachine"
	_ "github.com/ipfs/go-datastore"
	"github.com/qmuntal/stateless"
)

const (
	stateWait       = "Wait"
	statePrepairing = "Preparing"
	stateRun        = "Ringing"
	stateReward     = "Reward"
	stateFinish     = "Finish"
)

const (
	triggerPresenceReady    = "GamePresenceReady"
	triggerPrepairingDone   = "GamePrepairingDone"
	triggerPrepairingFailed = "GamePrepairingFailed"
	triggerRunTimeout       = "GameRunTimeout"
	triggerRunCombineAll    = "GameRunCombineAll"
	triggerRewardTimeout    = "GameRewardTimeout"
	triggerNoOne            = "GameNoOne"
)

type GameStateMachine struct {
	state *stateless.StateMachine
}

func (m *GameStateMachine) configure() {
	waitHandler := NewStateWait()
	m.state.Configure(stateWait).
		OnEntry(waitHandler.Enter).
		OnExit(waitHandler.Exit).
		Permit(triggerPresenceReady, statePrepairing).
		Permit(triggerNoOne, stateFinish)

	preparing := NewStatePrepairing()
	m.state.Configure(statePrepairing).
		OnEntry(preparing.Enter).
		OnExit(preparing.Exit).
		Permit(triggerPrepairingDone, stateRun).
		Permit(triggerPrepairingFailed, stateWait)

	run := NewStateRun()
	m.state.Configure(stateRun).
		OnEntry(run.Enter).
		OnExit(run.Exit).
		Permit(triggerRunTimeout, stateReward).
		Permit(triggerRunCombineAll, stateReward)

	reward := NewStateReward()
	m.state.Configure(stateReward).
		OnEntry(reward.Enter).
		OnExit(reward.Exit).
		Permit(triggerRewardTimeout, statePrepairing)

	m.state.ToGraph()
}

func (m *GameStateMachine) Fire(trigger stateless.Trigger, args ...interface{}) error {
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
