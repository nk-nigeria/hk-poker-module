package entity

import (
	"math"
	"math/rand"
	"time"

	pb "github.com/ciaolink-game-platform/cgp-chinese-poker-module/proto"
	"github.com/emirpasic/gods/maps/linkedhashmap"
	"github.com/heroiclabs/nakama-common/runtime"
)

const (
	TickRate                 = 5
	MaxEmptySec              = 60 * TickRate // 60s
	DelayBeforeRunGameSec    = 5 * TickRate  // 5s
	DelayBeforeRewardGameSec = 20 * TickRate // 30s
	DelayBeforeFinishGameSec = 10 * TickRate // 30s

	MinPresences = 2
)

type MatchLabel struct {
	Open              int32  `json:"open"`
	LastOpenValueNoti int32  `json:"-"` // using for check has noti new state of open
	Bet               int32  `json:"bet"`
	Code              string `json:"code"`
	Name              string `json:"name"`
	Password          string `json:"password"`
	MaxSize           int32  `json:"max_size"`
}

type MatchState struct {
	Random       *rand.Rand
	Label        *MatchLabel
	MinPresences int
	EmptyTicks   int

	// Currently connected users, or reserved spaces.
	Presences *linkedhashmap.Map
	// Number of users currently in the process of connecting to the match.
	JoinsInProgress int
	// Number of user currently dealt with game
	JoinInGame map[string]bool

	// Mark assignments to player user IDs.
	Cards map[string]*pb.ListCard
	// Mark assignments to player user IDs.
	OrganizeCards map[string]*pb.ListCard
	// Whose turn it currently is.

	gameState          pb.GameState
	CountDownReachTime time.Time
}

type CountDown struct {
	delayInit int64
	Tick      int64
	Sec       int64
	IsUpdate  bool
}

func NewCountDown(duration int64) CountDown {
	cd := CountDown{
		delayInit: duration,
		Tick:      duration,
	}
	cd.Sec = 0
	cd.IsUpdate = true
	return cd
}

func (cd *CountDown) doCountDown() {
	cd.Tick--
	if cd.Tick < 0 {
		return
	}
	v := math.Ceil(float64(cd.Tick) / float64(TickRate))
	if cd.Sec != int64(v) {
		cd.Sec = int64(v)
		cd.IsUpdate = true
	}
}

func (cd *CountDown) reset(sec int64) {
	cd.Tick = sec
	cd.Sec = cd.Tick / TickRate
	cd.IsUpdate = true
}

func NewMathState(label *MatchLabel) MatchState {
	m := MatchState{
		Random:       rand.New(rand.NewSource(time.Now().UnixNano())),
		Label:        label,
		MinPresences: MinPresences,
		Presences:    linkedhashmap.New(),
	}
	m.Label.LastOpenValueNoti = m.Label.Open
	return m
}

func PbGameStateString(gp pb.GameState) string {
	switch gp {
	case pb.GameState_GameStateWait:
		return "GameStateWait"
	case pb.GameState_GameStatePreparing:
		return "GameStatePreparing"
	case pb.GameState_GameStatePlay:
		return "GameStateRun"
	case pb.GameState_GameStateReward:
		return "GameStateReward"
	case pb.GameState_GameStateFinish:
		return "GameStateFinish"
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

func (s *MatchState) GetGameState() pb.GameState {
	return s.gameState
}

func (s *MatchState) GetNextGameState() pb.GameState {
	return (s.gameState + 1) % 7
}

//func (s *MatchState) SetGameState(gameState pb.GameState, logger runtime.Logger) pb.GameState {
//	if s.gameState != gameState {
//		logger.Info("Game state change %s -- > %s", s.gameState.String(), gameState.String())
//		s.gameState = gameState
//		// reset duration empty room
//		if s.gameState == pb.GameState_GameStateLobby {
//			s.EmptyTicks = 0
//		}
//		s.CountDown.reset(0)
//	}
//	return s.gameState
//}

//func (s *MatchState) ProcessEvent(gameEvent GameEvent, logger runtime.Logger, presences []runtime.Presence) *MatchState {
//	if gameEvent != MathLoop {
//		logger.Info("ProccessEvent %s, current gameState %s", gameEvent.String(), s.GetGameState().String())
//	}
//	switch s.gameState {
//	case pb.GameState_GameStateLobby:
//		s.handlerGameStateLobby(gameEvent, logger, presences)
//	case pb.GameState_GameStatePrepare:
//		s.handlerGamePrepare(gameEvent, logger, presences)
//	case pb.GameState_GameStateCountdown:
//		s.handlerGameCountDown(gameEvent, logger, presences)
//	case pb.GameState_GameStatePlay:
//		s.handlerGameRun(gameEvent, logger, presences)
//	case pb.GameState_GameStateReward:
//		s.handlerGameReward(gameEvent, logger, presences)
//	case pb.GameState_GameStateFinish:
//		s.handlerGameFinish(gameEvent, logger, presences)
//	}
//	return s
//}

func (s *MatchState) handlerGameStateLobby(gameEvent GameEvent, logger runtime.Logger, presences []runtime.Presence) *MatchState {
	//s.Label.Open = 1
	//if gameEvent == MatchJoin {
	//	s.addPresence(presences)
	//	s.SetGameState(pb.GameState_GameStatePrepare, logger)
	//	return s
	//}
	//
	//if gameEvent == MathLoop {
	//	if s.Presences.Size()+s.JoinsInProgress == 0 {
	//		s.EmptyTicks++
	//	}
	//}
	return s
}

func (s *MatchState) handlerGamePrepare(gameEvent GameEvent, logger runtime.Logger, presences []runtime.Presence) *MatchState {
	//s.Label.Open = 1
	//if gameEvent == MatchLeave {
	//	s.removePresence(presences)
	//	if s.Presences.Size() == 0 {
	//		s.SetGameState(pb.GameState_GameStateLobby, logger)
	//	}
	//	return s
	//}
	//if gameEvent == MatchJoin {
	//	s.addPresence(presences)
	//	if s.Presences.Size() >= MinPlayer {
	//		s.SetGameState(pb.GameState_GameStateCountdown, logger)
	//		s.CountDown.reset(DelayBeforeRunGameSec)
	//	}
	//	return s
	//}
	return s
}

func (s *MatchState) handlerGameCountDown(gameEvent GameEvent, logger runtime.Logger, presences []runtime.Presence) *MatchState {
	//if gameEvent == MatchLeave {
	//	s.removePresence(presences)
	//	return s
	//}
	//if gameEvent == MatchJoin {
	//	s.addPresence(presences)
	//	return s
	//}
	//
	//if gameEvent == MathLoop {
	//	s.CountDown.doCountDown()
	//	if s.CountDown.Tick < 0 {
	//		s.SetGameState(pb.GameState_GameStatePlay, logger)
	//		s.Cards = make(map[string]*pb.ListCard, 0) // clear map of list card
	//		s.CountDown.reset(DelayBeforeRewardGameSec)
	//		s.Label.Open = 0
	//	}
	//}
	return s
}

func (s *MatchState) handlerGameRun(gameEvent GameEvent, logger runtime.Logger, presences []runtime.Presence) *MatchState {
	//if gameEvent == MatchJoin {
	//	// todo handler add presence when game aldready run
	//	s.addPresence(presences)
	//
	//	return s
	//}
	//if gameEvent == MatchLeave {
	//	s.removePresence(presences)
	//	return s
	//}
	//if gameEvent == MathDone {
	//	s.SetGameState(pb.GameState_GameStateReward, logger)
	//	return s
	//}
	//if gameEvent == MathLoop {
	//	// todo punishment as looser
	//	if len(s.JoinInGame) <= 1 {
	//		s.SetGameState(pb.GameState_GameStateReward, logger)
	//		return s
	//	}
	//	s.CountDown.doCountDown()
	//	if s.CountDown.Tick < 0 {
	//		s.SetGameState(pb.GameState_GameStateReward, logger)
	//		s.CountDown.reset(DelayBeforeFinishGameSec)
	//		return s
	//	}
	//}
	// todo add param user commnad
	return s
}

func (s *MatchState) handlerGameReward(gameEvent GameEvent, logger runtime.Logger, presences []runtime.Presence) *MatchState {
	//if gameEvent == MatchJoin {
	//	s.addPresence(presences)
	//}
	//if gameEvent == MatchLeave {
	//	s.removePresence(presences)
	//}
	//// todo calc reward here
	//
	//s.CountDown.doCountDown()
	//if s.CountDown.Tick < 0 {
	//	s.SetGameState(pb.GameState_GameStateFinish, logger)
	//}
	return s
}

func (s *MatchState) handlerGameFinish(gameEvent GameEvent, logger runtime.Logger, presences []runtime.Presence) *MatchState {
	//s.Playing = false
	//if gameEvent == MatchJoin {
	//	s.addPresence(presences)
	//}
	//if gameEvent == MatchLeave {
	//	s.removePresence(presences)
	//}
	//if s.Presences.Size() >= MinPlayer {
	//	s.SetGameState(pb.GameState_GameStateCountdown, logger)
	//	s.CountDown.reset(DelayBeforeRunGameSec)
	//	return s
	//}
	//if s.Presences.Size() > 0 {
	//	s.SetGameState(pb.GameState_GameStatePrepare, logger)
	//	s.CountDown.reset(DelayBeforeRunGameSec)
	//	return s
	//}
	//s.SetGameState(pb.GameState_GameStateLobby, logger)
	return s
}

func (s *MatchState) handlerGameTerminate(gameEvent GameEvent, logger runtime.Logger, presences []runtime.Presence) *MatchState {
	// todo Game Terminate by server shutdown
	return s
}

func (s *MatchState) AddPresence(presences []runtime.Presence) {
	for _, presence := range presences {
		s.EmptyTicks = 0
		s.Presences.Put(presence.GetUserId(), presence)
		s.JoinsInProgress--
		if _, exist := s.Cards[presence.GetUserId()]; exist {
			s.JoinInGame[presence.GetUserId()] = true
		}
	}
}

func (s *MatchState) RemovePresence(presences []runtime.Presence) {
	for _, presence := range presences {
		s.Presences.Remove(presence.GetUserId())
		if _, exist := s.Cards[presence.GetUserId()]; exist {
			s.JoinInGame[presence.GetUserId()] = false
		}
	}
}

func (s *MatchState) SetUpCountDown(duration time.Duration) {
	s.CountDownReachTime = time.Now().Add(duration)
}

func (s *MatchState) GetRemainCountDown() int {
	currentTime := time.Now()
	difference := s.CountDownReachTime.Sub(currentTime)
	return int(difference.Seconds())
}
