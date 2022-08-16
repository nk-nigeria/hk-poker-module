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
	Open         int32  `json:"open"`
	Bet          int32  `json:"bet"`
	Code         string `json:"code"`
	Name         string `json:"name"`
	Password     string `json:"password"`
	MaxSize      int32  `json:"max_size"`
	MockCodeCard int32  `json:"mock_code_card"`
}

type MatchState struct {
	Random       *rand.Rand
	Label        *MatchLabel
	MinPresences int

	// Currently connected users, or reserved spaces.
	Presences        *linkedhashmap.Map
	PlayingPresences *linkedhashmap.Map
	LeavePresences   *linkedhashmap.Map
	// Number of users currently in the process of connecting to the match.
	JoinsInProgress     int
	PresencesNoInteract map[string]int

	// Mark assignments to player user IDs.
	Cards map[string]*pb.ListCard
	// Mark assignments to player user IDs.
	OrganizeCards map[string]*pb.ListCard
	// Whose turn it currently is.

	CountDownReachTime time.Time
	LastCountDown      int
}

func NewMathState(label *MatchLabel) MatchState {
	m := MatchState{
		Random:              rand.New(rand.NewSource(time.Now().UnixNano())),
		Label:               label,
		MinPresences:        MinPresences,
		Presences:           linkedhashmap.New(),
		PlayingPresences:    linkedhashmap.New(),
		LeavePresences:      linkedhashmap.New(),
		PresencesNoInteract: make(map[string]int, 0),
	}
	return m
}

func (s *MatchState) AddPresence(presences []runtime.Presence) {
	for _, presence := range presences {
		s.Presences.Put(presence.GetUserId(), presence)
		s.ResetUserNotInteract(presence.GetUserId())
	}
}

func (s *MatchState) RemovePresence(presences []runtime.Presence) {
	for _, presence := range presences {
		s.Presences.Remove(presence.GetUserId())
		delete(s.PresencesNoInteract, presence.GetUserId())
	}
}

func (s *MatchState) AddLeavePresence(presences []runtime.Presence) {
	for _, presence := range presences {
		s.LeavePresences.Put(presence.GetUserId(), presence)
	}
}

func (s *MatchState) ApplyLeavePresence() {
	s.LeavePresences.Each(func(key interface{}, value interface{}) {
		s.Presences.Remove(key)
		delete(s.PresencesNoInteract, key.(string))
	})

	s.LeavePresences = linkedhashmap.New()
}

func (s *MatchState) SetupMatchPresence() {
	s.PlayingPresences = linkedhashmap.New()
	s.Presences.Each(func(key interface{}, value interface{}) {
		s.PlayingPresences.Put(key, value)
		keyStr := key.(string)
		if val, exist := s.PresencesNoInteract[keyStr]; exist {
			s.PresencesNoInteract[keyStr] = val + 1
		} else {
			s.PresencesNoInteract[keyStr] = 1
		}
	})
}

func (s *MatchState) GetPresenceNotInteract(roundGame int) []runtime.Presence {
	listPresence := make([]runtime.Presence, 0)
	s.Presences.Each(func(key interface{}, value interface{}) {
		if roundGameNotInteract, exist := s.PresencesNoInteract[key.(string)]; exist && roundGameNotInteract >= roundGame {
			listPresence = append(listPresence, value.(runtime.Presence))
		}
	})
	return listPresence
}

func (s *MatchState) SetUpCountDown(duration time.Duration) {
	s.CountDownReachTime = time.Now().Add(duration)
	s.LastCountDown = -1
}

func (s *MatchState) GetRemainCountDown() int {
	currentTime := time.Now()
	difference := s.CountDownReachTime.Sub(currentTime)
	return int(difference.Seconds())
}

func (s *MatchState) SetLastCountDown(countDown int) {
	s.LastCountDown = countDown
}

func (s *MatchState) GetLastCountDown() int {
	return s.LastCountDown
}

func (s *MatchState) IsNeedNotifyCountDown() bool {
	return s.GetRemainCountDown() != s.LastCountDown || s.LastCountDown == -1
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
	// s.PresencesNoAction[userId] = 0
}

func (s *MatchState) GetPlayingCount() int {
	return s.PlayingPresences.Size()
}

func (s *MatchState) GetShowCardCount() int {
	return len(s.OrganizeCards)
}

func (s *MatchState) GetPresences() []runtime.Presence {
	presences := make([]runtime.Presence, 0)
	s.Presences.Each(func(key interface{}, value interface{}) {
		presences = append(presences, value.(runtime.Presence))
	})

	return presences
}

func (s *MatchState) GetPlayingPresences() []runtime.Presence {
	presences := make([]runtime.Presence, 0)
	s.PlayingPresences.Each(func(key interface{}, value interface{}) {
		presences = append(presences, value.(runtime.Presence))
	})

	return presences
}

func (s *MatchState) GetLeavePresences() []runtime.Presence {
	presences := make([]runtime.Presence, 0)
	s.LeavePresences.Each(func(key interface{}, value interface{}) {
		presences = append(presences, value.(runtime.Presence))
	})

	return presences
}

func (s *MatchState) ResetUserNotInteract(userId string) {
	s.PresencesNoInteract[userId] = 0
}
