package entity

import (
	"math/rand"
	"time"

	pb "github.com/ciaolink-game-platform/cgp-chinese-poker-module/proto"
	"github.com/emirpasic/gods/maps/linkedhashmap"
	"github.com/heroiclabs/nakama-common/runtime"
)

const (
	TickRate = 5

	MaxEmptySec = 60 * (1 / (TickRate * 0.1)) // 60s

	MinPlayer        = 2
	CountDownGameSec = 5 * (1 / (TickRate * 0.1)) // 5s
)

type MatchLabel struct {
	Open              int32  `json:"open"`
	LastOpenValueNoti int32  `json:"-"` // using for check has noti new state of open
	Bet               int32  `json:"bet"`
	Code              string `json:"code"`
	Name              string `json:"name"`
	Password          string `json:"password"`
}

type MatchState struct {
	Random     *rand.Rand
	Label      *MatchLabel
	EmptyTicks int

	// Currently connected users, or reserved spaces.
	Presences *linkedhashmap.Map
	// Number of users currently in the process of connecting to the match.
	JoinsInProgress int

	// True if there's a game currently in progress.
	Playing bool
	// Mark assignments to player user IDs.
	Cards map[string]*pb.ListCard
	// Mark assignments to player user IDs.
	OrganizeCards map[string]*pb.ListCard
	// Whose turn it currently is.
	Turn string
	// Ticks until they must submit their move.
	DeadlineRemainingTicks int64
	//// The winner positions.
	//winnerPositions []int32
	// Ticks until the next game starts, if applicable.
	NextGameRemainingTicks int64

	gameState          GamePresentState
	countDownEnterGame int64
}

func NewMathState(label *MatchLabel) MatchState {
	m := MatchState{
		Random:  rand.New(rand.NewSource(time.Now().UnixNano())),
		Label:   label,
		Playing: false,
		//presences: make(map[string]runtime.Presence, maxPlayer),
		Presences:          linkedhashmap.New(),
		gameState:          GameStateLobby,
		countDownEnterGame: 0,
	}
	m.Label.LastOpenValueNoti = m.Label.Open
	return m
}

type GamePresentState int

const (
	GameStateLobby GamePresentState = iota
	GameStatePrepare
	GameStateCountdown
	GameStateRun
	GameStateReward
	GameStateFinish
	GameStateEnd
)

func (gp GamePresentState) String() string {
	switch gp {
	case GameStateLobby:
		return "GameStateLobby"
	case GameStatePrepare:
		return "GameStatePrepare"
	case GameStateCountdown:
		return "GameStateCountdown"
	case GameStateRun:
		return "GameStateRun"
	case GameStateReward:
		return "GameStateReward"
	case GameStateFinish:
		return "GameStateFinish"
	case GameStateEnd:
		return "GameStateEnd"
	}
	return "unknow"
}

type GameEvent int

const (
	MatchJoin GameEvent = iota
	MatchJoinAttempt
	MatchLeave
	MathDone
	MathLoop
	MatchTerminate
)

func (ge GameEvent) String() string {
	switch ge {
	case MatchJoin:
		return "MatchJoin"
	case MatchJoinAttempt:
		return "MatchJoinAttempt"
	case MatchLeave:
		return "MatchLeave"
	case MathDone:
		return "MathDone"
	case MathLoop:
		return "MathLoop"
	case MatchTerminate:
		return "MatchTerminate"
	}
	return "unknow"
}

func (s *MatchState) GetGameState() GamePresentState {
	return s.gameState
}

func (s *MatchState) SetGameState(gameState GamePresentState, logger runtime.Logger) GamePresentState {
	if s.gameState != gameState {
		logger.Info("Game state change %s -- > %s", s.gameState.String(), gameState.String())
		s.gameState = gameState
		// reset duration empty room
		if s.gameState == GameStateLobby {
			s.EmptyTicks = 0
		}
	}
	return s.gameState
}
func (s *MatchState) ProcessEvent(gameEvent GameEvent, logger runtime.Logger, presences []runtime.Presence) *MatchState {
	if gameEvent != MathLoop {
		logger.Info("ProccessEvent %s, current gameState %s", gameEvent.String(), s.GetGameState().String())
	}
	switch s.gameState {
	case GameStateLobby:
		s.handlerGameStateLobby(gameEvent, logger, presences)
	case GameStatePrepare:
		s.handlerGamePrepare(gameEvent, logger, presences)
	case GameStateCountdown:
		s.handlerGameCountDown(gameEvent, logger, presences)
	case GameStateRun:
		s.handlerGameRun(gameEvent, logger, presences)
	case GameStateReward:
		s.handlerGameReward(gameEvent, logger, presences)
	case GameStateFinish:
		s.handlerGameFinish(gameEvent, logger, presences)
	}
	return s
}

func (s *MatchState) handlerGameStateLobby(gameEvent GameEvent, logger runtime.Logger, presences []runtime.Presence) *MatchState {
	s.Label.Open = 1
	if gameEvent == MatchJoin {
		s.addPresence(presences)
		s.SetGameState(GameStatePrepare, logger)
		return s
	}

	if gameEvent == MathLoop {
		if s.Presences.Size()+s.JoinsInProgress == 0 {
			s.EmptyTicks++
		}
	}
	return s
}

func (s *MatchState) handlerGamePrepare(gameEvent GameEvent, logger runtime.Logger, presences []runtime.Presence) *MatchState {
	s.Label.Open = 1
	if gameEvent == MatchLeave {
		s.removePresence(presences)
		if s.Presences.Size() == 0 {
			s.SetGameState(GameStateLobby, logger)
		}
		return s
	}
	if gameEvent == MatchJoin {
		s.addPresence(presences)
		if s.Presences.Size() >= MinPlayer {
			s.SetGameState(GameStateCountdown, logger)
			s.countDownEnterGame = CountDownGameSec
		}
		return s
	}
	return s
}

func (s *MatchState) handlerGameCountDown(gameEvent GameEvent, logger runtime.Logger, presences []runtime.Presence) *MatchState {
	if gameEvent == MatchLeave {
		s.removePresence(presences)
		if s.Presences.Size() == 0 {
			s.SetGameState(GameStateLobby, logger)
		}
		return s
	}
	if gameEvent == MatchJoin {
		s.addPresence(presences)
		return s
	}

	if gameEvent == MathLoop {
		s.countDownEnterGame--
		if s.countDownEnterGame <= 0 {
			s.SetGameState(GameStateRun, logger)
			s.Label.Open = 0
		}
	}
	return s
}

func (s *MatchState) handlerGameRun(gameEvent GameEvent, logger runtime.Logger, presences []runtime.Presence) *MatchState {
	if gameEvent == MatchJoin {
		// todo handler add presence when game aldready run
		return s
	}
	if gameEvent == MatchLeave {
		s.removePresence(presences)
		// todo punishment as looser
		if s.Presences.Size() == 1 {
			s.SetGameState(GameStateReward, logger)
		}
		return s
	}
	if gameEvent == MathDone {
		s.SetGameState(GameStateReward, logger)
		return s
	}
	if gameEvent == MathLoop {
	}
	// todo add param user commnad
	return s
}

func (s *MatchState) handlerGameReward(gameEvent GameEvent, logger runtime.Logger, presences []runtime.Presence) *MatchState {
	if gameEvent == MatchJoin {
		s.addPresence(presences)
	}
	if gameEvent == MatchLeave {
		s.removePresence(presences)
	}
	// todo calc reward here

	s.SetGameState(GameStateFinish, logger)
	return s
}

func (s *MatchState) handlerGameFinish(gameEvent GameEvent, logger runtime.Logger, presences []runtime.Presence) *MatchState {
	if gameEvent == MatchJoin {
		s.addPresence(presences)
	}
	if gameEvent == MatchLeave {
		s.removePresence(presences)
	}
	if s.Presences.Size() >= MinPlayer {
		s.SetGameState(GameStateCountdown, logger)
		return s
	}
	if s.Presences.Size() > 0 {
		s.SetGameState(GameStatePrepare, logger)
		return s
	}
	s.SetGameState(GameStateLobby, logger)
	return s
}

func (s *MatchState) handlerGameTerminate(gameEvent GameEvent, logger runtime.Logger, presences []runtime.Presence) *MatchState {
	// todo Game Terminate by server shutdown
	return s
}

func (s MatchState) addPresence(presences []runtime.Presence) {
	for _, presence := range presences {
		s.EmptyTicks = 0
		s.Presences.Put(presence.GetUserId(), presence)
		s.JoinsInProgress--
	}
}

func (s MatchState) removePresence(presences []runtime.Presence) {
	for _, presence := range presences {
		s.Presences.Remove(presence.GetUserId())
	}
}
