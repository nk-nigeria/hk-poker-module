package api

import (
	"sort"

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

	return nil
}

func (c *ChinesePokerGame) Organize(dispatcher runtime.MatchDispatcher, s *entity.MatchState, presence string, cards *pb.ListCard) error {
	s.OrganizeCards[presence] = cards
	return nil
}

func (c *ChinesePokerGame) Finish(dispatcher runtime.MatchDispatcher, s *entity.MatchState) {
	// Check every user
	// Check every hand
	// Calculate hand to point
}

func (c *ChinesePokerGame) IsHighcard(listCard []*pb.Card) ([]*pb.Card, bool) {
	// if _, valid := c.CheckStraightFlush(listCard); valid {
	// 	return nil, false
	// }
	if _, valid := c.CheckFourOfAKind(listCard); valid {
		return nil, false
	}
	if _, valid := c.CheckFullHouse(listCard); valid {
		return nil, false
	}
	if _, valid := c.CheckFlush(listCard); valid {
		return nil, false
	}
	if _, valid := c.CheckStraight(listCard); valid {
		return nil, false
	}
	if _, valid := c.CheckThreeOfAKind(listCard); valid {
		return nil, false
	}
	// if _, valid := c.CheckTwoPairs(listCard); valid {
	// 	return nil, false
	// }
	if _, valid := c.CheckPair(listCard); valid {
		return nil, false
	}
	return c.SortCard(listCard), true
}

// Thùng phá sảnh (en: Straight Flush)
// Năm lá bài cùng màu, đồng chất, cùng một chuỗi số
// Là Flush, có cùng chuỗi
func (c *ChinesePokerGame) CheckStraightFlush(listCard []*pb.Card) ([]*pb.Card, bool) {
	l := len(listCard)
	if l != 3 && l != 5 {
		return nil, false
	}
	listCard, valid := c.CheckFlush(listCard)
	if !valid {
		return nil, false
	}
	listCard, valid = c.CheckStraight(listCard)
	if !valid {
		return nil, false
	}
	return listCard, true
}

// Tứ quý (en: Four of a Kind)
// Bốn lá đồng số
func (c *ChinesePokerGame) CheckFourOfAKind(listCard []*pb.Card) ([]*pb.Card, bool) {
	l := len(listCard)
	if l != 3 && l != 5 {
		return nil, false
	}
	mapCardRank := make(map[int][]*pb.Card)
	for i := range listCard {
		card := listCard[i]
		rankPoint := entity.GetCardRankPoint(card.Rank)
		var list []*pb.Card
		if _, exist := mapCardRank[rankPoint]; exist {
			list = append(list, card)
		} else {
			list = make([]*pb.Card, 0)
			list = append(list, card)
		}
		mapCardRank[rankPoint] = list
	}
	if len(mapCardRank) != 2 {
		return nil, false
	}

	newListCard := make([]*pb.Card, 0, len(listCard))
	isFourOfAKind := false
	for _, list := range mapCardRank {
		if len(list) == 4 {
			isFourOfAKind = true
		}
		newListCard = append(newListCard, list...)
	}
	if isFourOfAKind {
		newListCard = c.SortCard(newListCard)
		return newListCard, true
	}
	return nil, false
}

// Cù lũ (en: Full House)
// Một bộ ba và một bộ đôi
// Bốn lá đồng số
func (c *ChinesePokerGame) CheckFullHouse(listCard []*pb.Card) ([]*pb.Card, bool) {
	l := len(listCard)
	if l != 3 && l != 5 {
		return nil, false
	}
	mapCardRank := make(map[int][]*pb.Card)
	for i := range listCard {
		card := listCard[i]
		rankPoint := entity.GetCardRankPoint(card.Rank)
		var list []*pb.Card
		if _, exist := mapCardRank[rankPoint]; exist {
			list = append(list, card)
		} else {
			list = make([]*pb.Card, 0)
			list = append(list, card)
		}
		mapCardRank[rankPoint] = list
	}
	if len(mapCardRank) != 2 {
		return nil, false
	}

	newListCard := make([]*pb.Card, 0, len(listCard))
	hasTriangle := false
	hasDouble := false

	for _, list := range mapCardRank {
		if len(list) == 3 {
			hasTriangle = true
			newListCard = append(list, newListCard...)
			continue
		}
		if len(list) == 2 {
			hasDouble = true
			newListCard = append(list, newListCard...)
			continue
		}
		newListCard = append(newListCard, list...)
	}
	if hasTriangle && hasDouble {
		newListCard = c.SortCard(newListCard)
		return newListCard, true
	}
	return nil, false
}

// Thùng (en: Flush)
// Năm lá bài cùng màu, đồng chất (nhưng không cùng một chuỗi số)
func (c *ChinesePokerGame) CheckFlush(listCard []*pb.Card) ([]*pb.Card, bool) {
	l := len(listCard)
	if l != 3 && l != 5 {
		return nil, false
	}
	listCard = c.SortCard(listCard)
	prevSuitPoint := entity.GetCardSuitPoint(listCard[0].GetSuit())
	for i := 1; i < len(listCard); i++ {
		card := listCard[i]
		suitPoint := entity.GetCardSuitPoint(card.GetSuit())
		if suitPoint != prevSuitPoint {
			return nil, false
		}
	}
	return listCard, true
}

// Sảnh (en: Straight)
// Năm lá bài trong một chuỗi số (nhưng không đồng chất)
func (c *ChinesePokerGame) CheckStraight(listCard []*pb.Card) ([]*pb.Card, bool) {
	l := len(listCard)
	if l != 3 && l != 5 {
		return nil, false
	}
	listCard = c.SortCard(listCard)
	prevRankPoint := entity.GetCardRankPoint(listCard[0].GetRank())
	for i := 1; i < len(listCard); i++ {
		card := listCard[i]
		rankPoint := entity.GetCardRankPoint(card.GetRank())
		prevRankPoint++
		if rankPoint != prevRankPoint {
			return nil, false
		}
	}
	return listCard, true
}

// Xám chi/Xám cô (en: Three of a Kind)
// Ba lá bài đồng số
func (c *ChinesePokerGame) CheckThreeOfAKind(listCard []*pb.Card) ([]*pb.Card, bool) {
	l := len(listCard)
	if l != 3 && l != 5 {
		return nil, false
	}
	mapCardRank := make(map[int][]*pb.Card)
	for i := range listCard {
		card := listCard[i]
		rankPoint := entity.GetCardRankPoint(card.Rank)
		var list []*pb.Card
		if _, exist := mapCardRank[rankPoint]; !exist {
			list = make([]*pb.Card, 0)
		}
		list = append(list, card)
		mapCardRank[rankPoint] = list
	}
	if len(mapCardRank) < 2 {
		return nil, false
	}

	newListCard := make([]*pb.Card, 0, len(listCard))
	hasTriangle := false

	for _, list := range mapCardRank {
		if len(list) == 3 {
			hasTriangle = true
			newListCard = append(list, newListCard...)
			continue
		}
		newListCard = append(newListCard, list...)
	}
	if hasTriangle {
		newListCard = c.SortCard(newListCard)
		return newListCard, true
	}
	return nil, false
}

// Thú (en: Two Pairs)
// Hai đôi
func (c *ChinesePokerGame) CheckTwoPairs(listCard []*pb.Card) ([]*pb.Card, bool) {
	l := len(listCard)
	if l != 3 && l != 5 {
		return nil, false
	}
	mapCardRank := make(map[int][]*pb.Card)
	for i := range listCard {
		card := listCard[i]
		rankPoint := entity.GetCardRankPoint(card.Rank)
		var list []*pb.Card
		if _, exist := mapCardRank[rankPoint]; !exist {
			list = make([]*pb.Card, 0)
		}
		list = append(list, card)
		mapCardRank[rankPoint] = list
	}
	if len(mapCardRank) < 3 {
		return nil, false
	}

	newListCard := make([]*pb.Card, 0, len(listCard))
	numPair := 0

	for _, list := range mapCardRank {
		if len(list) == 2 {
			numPair++
			newListCard = append(list, newListCard...)
			continue
		}
		newListCard = append(newListCard, list...)
	}
	if numPair == 2 {
		newListCard = c.SortCard(newListCard)
		return newListCard, true
	}
	return nil, false
}

// Đôi (en: Pair)
// Hai lá bài đồng số
func (c *ChinesePokerGame) CheckPair(listCard []*pb.Card) ([]*pb.Card, bool) {
	l := len(listCard)
	if l != 3 && l != 5 {
		return nil, false
	}
	mapCardRank := make(map[int][]*pb.Card)
	for i := range listCard {
		card := listCard[i]
		rankPoint := entity.GetCardRankPoint(card.Rank)
		var list []*pb.Card
		if _, exist := mapCardRank[rankPoint]; !exist {
			list = make([]*pb.Card, 0)
		}
		list = append(list, card)
		mapCardRank[rankPoint] = list
	}
	if len(mapCardRank) < 3 {
		return nil, false
	}

	newListCard := make([]*pb.Card, 0, len(listCard))
	numPair := 0

	for _, list := range mapCardRank {
		if len(list) == 2 {
			numPair++
			newListCard = append(list, newListCard...)
			continue
		}
		newListCard = append(newListCard, list...)
	}
	if numPair > 0 {
		newListCard = c.SortCard(newListCard)
		return newListCard, true
	}
	return nil, false
}

// sort card increament by rank, equal rank will check suit
func (c *ChinesePokerGame) SortCard(listCard []*pb.Card) []*pb.Card {
	sort.Slice(listCard, func(a, b int) bool {
		cardA := listCard[a]
		cardB := listCard[b]
		rankPointA := entity.GetCardRankPoint(cardA.GetRank())
		rankPointB := entity.GetCardRankPoint(cardB.GetRank())
		if rankPointA > rankPointB {
			return true
		}
		if rankPointA < rankPointB {
			return false
		}
		suitPointA := entity.GetCardSuitPoint(cardA.GetSuit())
		suitPointB := entity.GetCardSuitPoint(cardB.GetSuit())
		return suitPointA < suitPointB
	})
	return listCard
}
