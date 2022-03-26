package chinese_poker

import (
	"github.com/ciaolink-game-platform/cgp-chinese-poker-module/entity"
	"github.com/ciaolink-game-platform/cgp-chinese-poker-module/pkg/log"
	pb "github.com/ciaolink-game-platform/cgp-chinese-poker-module/proto"
	"github.com/mxschmitt/golang-combinations"
)

const MaxPresenceCard = 13

type Engine struct {
	deck *entity.Deck
}

func NewChinesePokerEngine() UseCase {
	return &Engine{}
}

func (c *Engine) NewGame(s *entity.MatchState) error {
	s.JoinInGame = make(map[string]bool)
	s.Cards = make(map[string]*pb.ListCard)
	s.OrganizeCards = make(map[string]*pb.ListCard)

	return nil
}

func (c *Engine) Deal(s *entity.MatchState) error {
	c.deck = entity.NewDeck()
	c.deck.Shuffle()

	// loop on userid in match
	for _, k := range s.PlayingPresences.Keys() {
		userId := k.(string)
		cards, err := c.deck.Deal(MaxPresenceCard)
		if err == nil {
			s.Cards[userId] = cards
			s.JoinInGame[userId] = true
		} else {
			return err
		}
	}

	return nil
}

func (c *Engine) Organize(s *entity.MatchState, presence string, cards *pb.ListCard) error {
	s.UpdateShowCard(presence, cards)
	return nil
}

func (c *Engine) Combine(s *entity.MatchState, presence string) error {
	s.RemoveShowCard(presence)
	return nil
}

func (c *Engine) Finish(s *entity.MatchState) *pb.UpdateFinish {
	// Check every user
	updateFinish := pb.UpdateFinish{}
	presenceCount := s.PlayingPresences.Size()
	ctx := NewCompareContext(presenceCount)

	log.GetLogger().Info("Finish presence %v, size %v", s.PlayingPresences, presenceCount)

	// prepare for compare data
	userIds := make([]string, 0, presenceCount)
	hands := make(map[string]*Hand)
	results := make(map[string]*pb.ComparisonResult)
	for _, val := range s.PlayingPresences.Keys() {
		uid := val.(string)
		userIds = append(userIds, uid)

		cards := s.OrganizeCards[uid]
		var hand *Hand
		var err error
		hand, err = NewHandFromPb(cards)
		hand.SetOwner(uid)
		if err != nil {
			continue
		}

		hand.calculatePoint()
		hands[uid] = hand

		result := &pb.ComparisonResult{
			UserId:      uid,
			PointResult: hand.GetPointResult(),
			ScoreResult: &pb.ScoreResult{},
		}

		results[uid] = result

		updateFinish.Results = append(updateFinish.Results, result)

		log.GetLogger().Info("prepare for %s, hand %v, result %v", uid, hand, result)
	}

	pairs := combinations.Combinations(userIds, 2)
	log.GetLogger().Info("combination %v of %v", pairs, len(userIds))
	for _, pair := range pairs {
		uid1 := pair[0]
		uid2 := pair[1]
		log.GetLogger().Info("compare %v with %v", pair[0], pair[1])

		// calculate natural point, normal point, hand bonus case
		rc := CompareHand(ctx, hands[uid1], hands[uid2])
		ProcessCompareResult(ctx, results[uid1], rc.r1)
		ProcessCompareResult(ctx, results[uid2], rc.r2)

		updateFinish.Bonuses = append(updateFinish.Bonuses, rc.bonuses...)
	}

	ProcessCompareBonusResult(ctx, updateFinish.Results, &updateFinish.Bonuses)
	CalcTotalFactor(updateFinish.Results)

	return &updateFinish
}
