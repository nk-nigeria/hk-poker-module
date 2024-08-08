package bin_list_card

import "github.com/nakamaFramework/cgp-chinese-poker-module/entity"

type CombineType = int

var (
	CombinePair      CombineType = 1
	CombineThree     CombineType = 2
	CombineFour      CombineType = 3
	CombineStraight  CombineType = 4
	CombineFullHouse CombineType = 5
	CombineFlush     CombineType = 6
	CombineFullColor CombineType = 7
)

type ChinesePokerBinList interface {
	GetChain(b *entity.BinListCard, comb CombineType) (uint, entity.ListCard)
}
