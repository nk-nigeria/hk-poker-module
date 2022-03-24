package chinese_poker

import (
	"github.com/ciaolink-game-platform/cgp-chinese-poker-module/entity"
	pb "github.com/ciaolink-game-platform/cgp-chinese-poker-module/proto"
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
	ctx := NewCompareContext(s.PlayingPresences.Size())

	for _, uid1 := range s.PlayingPresences.Keys() {
		userID1 := uid1.(string)
		cards1 := s.OrganizeCards[userID1]
		hand1, err := NewHand(cards1)
		if err != nil {
			continue
		}

		hand1.calculatePoint()

		result := &pb.ComparisonResult{
			UserId:      userID1,
			PointResult: hand1.GetPointResult(),
			ScoreResult: &pb.ScoreResult{},
		}

		for _, uid2 := range s.PlayingPresences.Keys() {
			userID2 := uid2.(string)
			if userID1 == userID2 {
				continue
			}
			cards2 := s.OrganizeCards[userID2]
			hand2, err := NewHand(cards2)
			if err != nil {
				continue
			}

			// calculate natural point, normal point, hand bonus case
			rc := CompareHand(ctx, hand1, hand2)
			ProcessCompareResult(ctx, result, rc)
		}

		updateFinish.Results = append(updateFinish.Results, result)
	}

	ProcessCompareBonusResult(ctx, updateFinish.Results)
	CalcTotalFactor(updateFinish.Results)

	return &updateFinish
}
