package entity

import (
	"math/rand"
	"time"

	pb "github.com/ciaolink-game-platform/cgp-chinese-poker-module/proto"
	"github.com/emirpasic/gods/maps/linkedhashmap"
	"github.com/heroiclabs/nakama-common/runtime"
)

const (
	MinPresences = 2
	MaxPresences = 4
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
	Presences        *linkedhashmap.Map
	PlayingPresences *linkedhashmap.Map
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

func (s *MatchState) SetupPlayingPresence() {
	s.PlayingPresences = linkedhashmap.New()
	s.Presences.Each(func(key interface{}, value interface{}) {
		s.PlayingPresences.Put(key, value)
	})
}

func (s *MatchState) SetUpCountDown(duration time.Duration) {
	s.CountDownReachTime = time.Now().Add(duration)
}

func (s *MatchState) GetRemainCountDown() int {
	currentTime := time.Now()
	difference := s.CountDownReachTime.Sub(currentTime)
	return int(difference.Seconds())
}

func (s *MatchState) GetPresenceSize() int {
	return s.Presences.Size()
}

func (s *MatchState) GetPlayingPresenceSize() int {
	return s.PlayingPresences.Size()
}

func (s *MatchState) IsReadyToPlay() bool {
	return s.Presences.Size() >= s.MinPresences
}

func (s *MatchState) UpdateShowCard(userId string, cards *pb.ListCard) {
	s.OrganizeCards[userId] = cards
}

func (s *MatchState) RemoveShowCard(userId string) {
	delete(s.OrganizeCards, userId)
}

func (s *MatchState) GetPlayingCount() int {
	return s.PlayingPresences.Size()
}

func (s *MatchState) GetShowCardCount() int {
	return len(s.OrganizeCards)
}
