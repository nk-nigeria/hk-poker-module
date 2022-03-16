package api

import (
	"github.com/ciaolink-game-platform/cgp-chinese-poker-module/entity"
	_ "github.com/filecoin-project/go-statemachine"
	"github.com/heroiclabs/nakama-common/runtime"
	_ "github.com/ipfs/go-datastore"
	"github.com/qmuntal/stateless"
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

type GameStateMachine struct {
	state *stateless.StateMachine
}

func (m *GameStateMachine) configure() {
	wait := NewStateWait(m.state.Fire)
	m.state.Configure(stateWait).
		OnEntry(wait.Enter).
		OnExit(wait.Exit).
		InternalTransition(triggerProcessWait, wait.Process).
		Permit(triggerPresenceReady, statePreparing).
		Permit(triggerNoOne, stateFinish)

	preparing := NewStatePreparing(m.state.Fire)
	m.state.Configure(statePreparing).
		OnEntry(preparing.Enter).
		OnExit(preparing.Exit).
		InternalTransition(triggerProcessPreparing, preparing.Process).
		Permit(triggerPreparingDone, statePlay).
		Permit(triggerPreparingFailed, stateWait)

	play := NewStatePlay()
	m.state.Configure(statePlay).
		OnEntry(play.Enter).
		OnExit(play.Exit).
		InternalTransition(triggerProcessPlay, play.Process).
		Permit(triggerPlayTimeout, stateReward).
		Permit(triggerPlayCombineAll, stateReward)

	reward := NewStateReward()
	m.state.Configure(stateReward).
		OnEntry(reward.Enter).
		OnExit(reward.Exit).
		InternalTransition(triggerProcessReward, reward.Process).
		Permit(triggerRewardTimeout, statePreparing)

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
	case statePreparing:
		trigger = triggerProcessPreparing
	case statePlay:
		trigger = triggerProcessPlay
	case stateReward:
		trigger = triggerProcessReward
	}
	return m.state.Fire(trigger, args...)
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

type FireFn func(trigger stateless.Trigger, args ...interface{}) error
type StateBase struct {
	fireFn FireFn
}

func GetState(args ...interface{}) *entity.MatchState {
	return args[0].(*entity.MatchState)
}

func GetLogger(args ...interface{}) runtime.Logger {
	return args[1].(runtime.Logger)
}

func GetDispatcher(args ...interface{}) runtime.MatchDispatcher {
	return args[2].(runtime.MatchDispatcher)
}

func GetProcessor(args ...interface{}) *MatchProcessor {
	return args[3].(*MatchProcessor)
}
