package api

import (
	"github.com/ciaolink-game-platform/cgp-chinese-poker-module/entity"
	pb "github.com/ciaolink-game-platform/cgp-chinese-poker-module/proto"
	"github.com/heroiclabs/nakama-common/runtime"
)

const MaxPresenceCard = 13

type ChinesePokerGame struct {
	deck *entity.Deck
}

func NewProcessor() *ChinesePokerGame {
	return &ChinesePokerGame{}
}

func (c *ChinesePokerGame) NewGame(s *entity.MatchState) error {
	s.OrganizeCards = make(map[string]*pb.ListCard)

	return nil
}

func (c *ChinesePokerGame) Deal(s *entity.MatchState) error {
	c.deck = entity.NewDeck()
	c.deck.Shuffle()

	s.JoinInGame = make(map[string]bool)
	s.Cards = make(map[string]*pb.ListCard)
	// loop on userid in match
	for _, k := range s.Presences.Keys() {
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

func (c *ChinesePokerGame) Organize(dispatcher runtime.MatchDispatcher, s *entity.MatchState, presence string, cards *pb.ListCard) error {
	s.OrganizeCards[presence] = cards
	return nil
}

func (c *ChinesePokerGame) Finish(dispatcher runtime.MatchDispatcher, s *entity.MatchState) *pb.UpdateFinish {
	// Check every user
	updateFinish := pb.UpdateFinish{}
	for _, uid1 := range s.Presences.Keys() {
		userID1 := uid1.(string)
		result := pb.ComparisonResult{
			UserId: userID1,
		}
		cards1 := s.Cards[userID1]
		hand1, err := NewHand(cards1)
		if err != nil {
			continue
		}
		for _, uid2 := range s.Presences.Keys() {
			userID2 := uid2.(string)
			if userID1 == userID2 {
				continue
			}
			cards2 := s.Cards[userID2]
			hand2, err := NewHand(cards2)
			if err != nil {
				continue
			}
			r := CompareHand(hand1, hand2)
			result.FrontFactor += r.FrontFactor
			result.MiddleFactor += r.MiddleFactor
			result.BackFactor += r.BackFactor
			result.FrontBonus += r.FrontBonus
			result.MiddleBonus += r.MiddleBonus
			result.BackBonus += r.BackBonus
		}
		updateFinish.Results = append(updateFinish.Results, &result)
	}
	return &updateFinish
	// Check every hand
	// Calculate hand to point
}
