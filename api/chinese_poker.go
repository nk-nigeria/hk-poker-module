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
	c.deck = entity.NewDeck()
	c.deck.Shuffle()

	s.Cards = make(map[string]*pb.ListCard)
	// loop on userid in match
	for _, k := range s.Presences.Keys() {
		userId := k.(string)
		cards, err := c.deck.Deal(MaxPresenceCard)
		if err == nil {
			s.Cards[userId] = cards
		} else {
			return err
		}
	}

	s.OrganizeCards = make(map[string]*pb.ListCard)

	return nil
}

func (c *ChinesePokerGame) Organize(dispatcher runtime.MatchDispatcher, s *entity.MatchState, presence string, cards *pb.ListCard) error {
	s.OrganizeCards[presence] = cards
	return nil
}

func (c *ChinesePokerGame) FinishGame(dispatcher runtime.MatchDispatcher, s *entity.MatchState) {

}
