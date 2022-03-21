package entity

import (
	pb "github.com/ciaolink-game-platform/cgp-chinese-poker-module/proto"
)

const (
	ModuleName = "chinese-poker"
)

var (
	mapHandRankingPoint map[pb.HandRanking]int
	mapWinFactorBonus   map[int64]int
)

func Min(x int, y int) int {
	if x < y {
		return x
	}
	return y
}

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

	// clear win
	// mapWinFactorBonus[pb.WinType_WIN_TYPE_LOSE_ALL_PLAYER] = -3
	// mapWinFactorBonus[pb.WinType_WIN_TYPE_LOSE_ALL_PLAYER] = -3
	mapWinFactorBonus = make(map[int64]int)
	mapWinFactorBonus[WIN_TYPE_WIN_CLEAN_DRAGON] = 108
	mapWinFactorBonus[WIN_TYPE_WIN_DRAGON] = 36
	mapWinFactorBonus[WIN_TYPE_WIN_FULL_COLORED] = 30
	mapWinFactorBonus[WIN_TYPE_WIN_FIVE_PAIR_THREE_OF_A_KIND] = 20
	mapWinFactorBonus[WIN_TYPE_WIN_SIX_AND_A_HALF_PAIRS] = 20
	mapWinFactorBonus[WIN_TYPE_WIN_THREE_STRAIGHT_FLUSH] = 20
	mapWinFactorBonus[WIN_TYPE_WIN_THREE_FLUSH] = 18
	mapWinFactorBonus[WIN_TYPE_WIN_THREE_STRAIGHT] = 16

	// bonus win
	mapWinFactorBonus[WIN_TYPE_WIN_FRONT_THREE_OF_A_KIND] = 3
	mapWinFactorBonus[WIN_TYPE_WIN_MID_FULL_HOUSE] = 2
	mapWinFactorBonus[WIN_TYPE_WIN_BACK_FOUR_OF_A_KIND] = 4
	mapWinFactorBonus[WIN_TYPE_WIN_MID_FOUR_OF_A_KIND] = 8
	mapWinFactorBonus[WIN_TYPE_WIN_BACK_STRAIGHT_FLUSH] = 5
	mapWinFactorBonus[WIN_TYPE_WIN_MID_STRAIGHT_FLUSH] = 10

	mapWinFactorBonus[WIN_TYPE_WIN_SCOOP] = 3
}

func GetHandRankingPoint(handRanking pb.HandRanking) int {
	return mapHandRankingPoint[handRanking]
}

func GetWinFactorBonus(winType int64) int {
	return mapWinFactorBonus[winType]
}
func Max(a, b int) int {
	if a < b {
		return b
	}
	return a
}
