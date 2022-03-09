package entity

import (
	pb "github.com/ciaolink-game-platform/cgp-chinese-poker-module/proto"
)

const (
	ModuleName = "chinese-poker"
)

var mapHandRankingPoint map[pb.HandRanking]int
var mapCardRankPoint map[pb.CardRank]int
var mapCardSuitPoint map[pb.CardSuit]int

func init() {
	// init hand ranking point
	mapHandRankingPoint = make(map[pb.HandRanking]int)
	mapHandRankingPoint[pb.HandRanking_HighCard] = 1
	mapHandRankingPoint[pb.HandRanking_Pair] = 2
	mapHandRankingPoint[pb.HandRanking_TwoPairs] = 3
	mapHandRankingPoint[pb.HandRanking_ThreeOfAKind] = 4
	mapHandRankingPoint[pb.HandRanking_Straight] = 5
	mapHandRankingPoint[pb.HandRanking_Flush] = 6
	mapHandRankingPoint[pb.HandRanking_FullHouse] = 7
	mapHandRankingPoint[pb.HandRanking_FourOfAKind] = 8
	mapHandRankingPoint[pb.HandRanking_StraightFlush] = 9

	mapCardRankPoint = make(map[pb.CardRank]int)
	mapCardRankPoint[pb.CardRank_RANK_2] = 1
	mapCardRankPoint[pb.CardRank_RANK_3] = 2
	mapCardRankPoint[pb.CardRank_RANK_4] = 3
	mapCardRankPoint[pb.CardRank_RANK_5] = 4
	mapCardRankPoint[pb.CardRank_RANK_6] = 5
	mapCardRankPoint[pb.CardRank_RANK_7] = 6
	mapCardRankPoint[pb.CardRank_RANK_8] = 7
	mapCardRankPoint[pb.CardRank_RANK_9] = 8
	mapCardRankPoint[pb.CardRank_RANK_10] = 8
	mapCardRankPoint[pb.CardRank_RANK_J] = 10
	mapCardRankPoint[pb.CardRank_RANK_Q] = 11
	mapCardRankPoint[pb.CardRank_RANK_K] = 12
	mapCardRankPoint[pb.CardRank_RANK_A] = 13

	mapCardSuitPoint = make(map[pb.CardSuit]int)
	mapCardSuitPoint[pb.CardSuit_SUIT_UNSPECIFIED] = 0
	mapCardSuitPoint[pb.CardSuit_SUIT_SPADES] = 1
	mapCardSuitPoint[pb.CardSuit_SUIT_CLUBS] = 2
	mapCardSuitPoint[pb.CardSuit_SUIT_DIAMONDS] = 3
	mapCardSuitPoint[pb.CardSuit_SUIT_HEARTS] = 4

}

func GetHandRankingPoint(handRanking pb.HandRanking) int {
	return mapHandRankingPoint[handRanking]
}

func GetCardRankPoint(cardRank pb.CardRank) int {
	return mapCardRankPoint[cardRank]
}

func GetCardSuitPoint(cardSuit pb.CardSuit) int {
	return mapCardSuitPoint[cardSuit]
}
