package entity

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	"github.com/ciaolink-game-platform/cgp-common/bot"
	"github.com/ciaolink-game-platform/cgp-common/define"
	"github.com/ciaolink-game-platform/cgp-common/lib"

	pb "github.com/ciaolink-game-platform/cgp-common/proto"
	"github.com/emirpasic/gods/maps/linkedhashmap"
	"github.com/heroiclabs/nakama-common/runtime"
)

var BotLoader = bot.NewBotLoader(nil, "", 0)
var GameStateDuration = lib.GetGameStateDurationByGameCode(define.ChinesePoker)

const (
	TickRate = 2

	MinPresences  = 2
	MaxPresences  = 4
	MinJpTreasure = 250 * 1000 * 1000
)

var RatioJpByVip = func(vipLv int64) float64 {
	if vipLv <= 1 {
		return 1
	}
	if vipLv <= 4 {
		return 2
	}
	return 2.5
}

type MatchState struct {
	Random       *rand.Rand
	Label        *pb.Match
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
	// GameState          pb.GameState
	// save balance result in state reward
	// using for send noti to presence join in state reward
	balanceResult    *pb.BalanceResult
	jackpotTreasure  *pb.Jackpot
	messages         []runtime.MatchData
	Bots             []*bot.BotPresence
	TurnOfBots       []*bot.BotPresence
	LastMoveCardUnix map[string]int64
	DelayForDeclare  *lib.TickCountDown
	MatchCount       int
	CardEvent        map[string]pb.CardEvent
}

func NewMathState(label *pb.Match) MatchState {
	m := MatchState{
		Random:              rand.New(rand.NewSource(time.Now().UnixNano())),
		Label:               label,
		MinPresences:        MinPresences,
		Presences:           linkedhashmap.New(),
		PlayingPresences:    linkedhashmap.New(),
		LeavePresences:      linkedhashmap.New(),
		PresencesNoInteract: make(map[string]int, 0),
		balanceResult:       nil,
		Bots:                make([]*bot.BotPresence, 0),
		LastMoveCardUnix:    map[string]int64{},
		DelayForDeclare:     &lib.TickCountDown{},
		CardEvent:           map[string]pb.CardEvent{},
	}
	if bots, err := BotLoader.GetFreeBot(int(label.GetNumBot())); err == nil {
		m.Bots = bots
	} else {
		fmt.Printf("\r\n load bot failed %s  \r\n", err.Error())
	}
	// Automatically add bot players
	// if label.GetNumBot() > 0 {
	// 	m.Bots = append(m.Bots, bot.NewBotPresences(int(label.GetNumBot()))...)
	// }
	for _, bot := range m.Bots {
		m.Presences.Put(bot.GetUserId(), bot)
	}
	return m
}

func (s *MatchState) Init() {
	s.Cards = make(map[string]*pb.ListCard)
	s.OrganizeCards = make(map[string]*pb.ListCard)
	s.LastMoveCardUnix = make(map[string]int64)
	playtimeout := GameStateDuration[pb.GameState_GameStatePlay].Seconds()
	for idx, v := range s.Bots {
		opt := bot.TurnOpt{MinTick: 5 * TickRate, MaxTick: int(playtimeout-playtimeout/3) * TickRate, MaxOccur: 1}
		v.InitTurnWithOption(opt, func() {
			x := s.Bots[idx]
			s.BotTurn(x)
		})
	}
	s.CardEvent = make(map[string]pb.CardEvent)
	for _, precense := range s.GetPresences() {
		s.CardEvent[precense.GetUserId()] = pb.CardEvent_COMBINE
	}
}

func (s *MatchState) GetBalanceResult() *pb.BalanceResult {
	return s.balanceResult
}
func (s *MatchState) SetBalanceResult(u *pb.BalanceResult) {
	s.balanceResult = u
}

func (s *MatchState) ResetBalanceResult() {
	s.SetBalanceResult(nil)
}

func (s *MatchState) SetJackpotTreasure(jp *pb.Jackpot) {
	s.jackpotTreasure = &pb.Jackpot{
		GameCode: jp.GetGameCode(),
		Chips:    jp.GetChips(),
	}
}
func (s *MatchState) GetJackpotTreasure() *pb.Jackpot {
	return s.jackpotTreasure
}

func (s *MatchState) AddPresence(ctx context.Context, db *sql.DB, presences []runtime.Presence) {
	for _, presence := range presences {
		m := NewMyPrecense(ctx, db, presence)
		s.Presences.Put(presence.GetUserId(), m)
		s.ResetUserNotInteract(presence.GetUserId())
	}
}

func (s *MatchState) RemovePresence(presences ...runtime.Presence) {
	for _, presence := range presences {
		s.Presences.Remove(presence.GetUserId())
		delete(s.PresencesNoInteract, presence.GetUserId())
		s.PlayingPresences.Remove(presence.GetUserId())
	}
}

func (s *MatchState) AddLeavePresence(presences ...runtime.Presence) {
	for _, presence := range presences {
		s.LeavePresences.Put(presence.GetUserId(), presence)
	}
}

func (s *MatchState) RemoveLeavePresence(userId string) {
	s.LeavePresences.Remove(userId)
}

func (s *MatchState) ApplyLeavePresence() {
	s.LeavePresences.Each(func(key interface{}, value interface{}) {
		// s.Presences.Remove(key)
		// delete(s.PresencesNoInteract, key.(string))
		s.RemovePresence(value.(runtime.Presence))
	})
	s.LeavePresences = linkedhashmap.New()
}

func (s *MatchState) SetupMatchPresence() {
	s.PlayingPresences = linkedhashmap.New()
	presences := make([]runtime.Presence, 0, s.Presences.Size())
	s.Presences.Each(func(key interface{}, value interface{}) {
		presences = append(presences, value.(runtime.Presence))
	})
	s.AddPlayingPresences(presences...)
}

func (s *MatchState) AddPlayingPresences(presences ...runtime.Presence) {
	for _, p := range presences {
		s.PlayingPresences.Put(p.GetUserId(), p)
		keyStr := p.GetUserId()
		if val, exist := s.PresencesNoInteract[keyStr]; exist {
			s.PresencesNoInteract[keyStr] = val + 1
		} else {
			s.PresencesNoInteract[keyStr] = 1
		}
	}
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

func (s *MatchState) GetPlayingNotBotCount() int {
	num := s.PlayingPresences.Size()
	s.PlayingPresences.Each(func(key, value interface{}) {
		p := value.(runtime.Presence)
		if BotLoader.IsBot(p.GetUserId()) {
			num--
		}
	})
	return num
}

func (s *MatchState) GetPrecenseNotBotCount() int {
	num := s.Presences.Size()
	s.Presences.Each(func(key, value interface{}) {
		p := value.(runtime.Presence)
		if BotLoader.IsBot(p.GetUserId()) {
			num--
		}
	})
	return num
}

func (s *MatchState) GetPrecenseBotCount() int {
	num := 0
	s.Presences.Each(func(key, value interface{}) {
		p := value.(runtime.Presence)
		if BotLoader.IsBot(p.GetUserId()) {
			num++
		}
	})
	return num
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

func (s *MatchState) GetPresence(userID string) runtime.Presence {
	_, value := s.Presences.Find(func(key, value interface{}) bool {
		if key == userID {
			return true
		}
		return false
	})
	if value == nil {
		return nil
	}
	return value.(runtime.Presence)
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

func (s *MatchState) BotTurn(v *bot.BotPresence) error {
	s.TurnOfBots = append(s.TurnOfBots, v)
	// find card of bot
	// var botCard *pb.ListCard
	// for userId, card := range s.Cards {
	// 	if userId == v.GetUserId() {
	// 		botCard = card
	// 		break
	// 	}
	// }
	// //
	// buf, err := defaultMarshaler.Marshal(&pb.Organize{
	// 	Cards: &pb.ListCard{
	// 		Cards: botCard.Cards,
	// 	},
	// })
	// if err != nil {
	// 	return err
	// }
	// reqs := []pb.OpCodeRequest{
	// 	pb.OpCodeRequest_OPCODE_REQUEST_DECLARE_CARDS,
	// }
	// for _, req := range reqs {
	// 	data := bot.NewBotMatchData(
	// 		req, buf, v,
	// 	)
	// 	s.messages = append(s.messages, data)
	// }
	return nil
}

func (s *MatchState) FakeDeclardCards(v *bot.BotPresence, cards []*pb.Card) {
	buf, err := defaultMarshaler.Marshal(&pb.Organize{
		Cards: &pb.ListCard{
			Cards: cards,
		},
	})
	if err != nil {
		return
	}
	reqs := []pb.OpCodeRequest{
		pb.OpCodeRequest_OPCODE_REQUEST_DECLARE_CARDS,
	}
	for _, req := range reqs {
		data := bot.NewBotMatchData(
			req, buf, v,
		)
		s.messages = append(s.messages, data)
	}
}

func (s *MatchState) BotLoop() {
	s.TurnOfBots = nil
	for _, v := range s.Bots {
		v.Loop()
	}
}

func (s *MatchState) Messages() []runtime.MatchData {
	msgs := s.messages
	s.messages = make([]runtime.MatchData, 0)
	return msgs
}

//	func (s *MatchState) AutoSortCard(cards []*pb.Card) []*pb.Card {
//		// 0:2
//		// 3:8
//		// 8:14
//		ml := NewBinListCards(NewListCard(cards)).ToList()
//		return nil
//	}
func (s *MatchState) UpdateLabel() {
	s.Label.Size = int32(s.GetPresenceSize())
	s.Label.Profiles = make([]*pb.SimpleProfile, 0)
	for _, precense := range s.GetPresences() {
		player := NewPlayer(precense)
		s.Label.Profiles = append(s.Label.Profiles, &pb.SimpleProfile{
			UserId:   precense.GetUserId(),
			UserName: precense.GetUsername(),
			UserSid:  player.Sid,
		})
	}
}

func (s *MatchState) GetGameState() pb.GameState {
	return s.Label.GameState
}
