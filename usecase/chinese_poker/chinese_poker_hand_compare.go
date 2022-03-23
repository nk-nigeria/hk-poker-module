package chinese_poker

import (
	"context"
	pb "github.com/ciaolink-game-platform/cgp-chinese-poker-module/proto"
)

var (
	mapNaturalPoint = map[pb.HandBonusType]int{
		pb.HandBonusType_BonusNaturalCleanDragon:    15,
		pb.HandBonusType_BonusNaturalDragon:         14,
		pb.HandBonusType_BonusNaturalFullColors:     13,
		pb.HandBonusType_BonusNaturalSixPairs:       12,
		pb.HandBonusType_BonusNaturalThreeOfFlushes: 11,
		pb.HandBonusType_BonusNaturalThreeStraights: 10,
	}

	mapBonusPoint = map[pb.HandBonusType]int{
		pb.HandBonusType_BonusThreeOfAKindFrontHand: 3,
		pb.HandBonusType_BonusFullHouseMidHand:      2,
		pb.HandBonusType_BonusFourOfAKindMidHand:    8,
		pb.HandBonusType_BonusStraightFlushMidHand:  10,
		pb.HandBonusType_BonusFourOfAKindBackHand:   4,
		pb.HandBonusType_BonusStraightFlushBackHand: 5,

		pb.HandBonusType_Scoop: 6,
	}
)

func rankingTypeToBonusType(ranking pb.HandRanking) pb.HandBonusType {
	switch ranking {
	case pb.HandRanking_NaturalCleanDragon:
		return pb.HandBonusType_BonusNaturalCleanDragon
	case pb.HandRanking_NaturalDragon:
		return pb.HandBonusType_BonusNaturalDragon
	case pb.HandRanking_NaturalFullColors:
		return pb.HandBonusType_BonusNaturalFullColors
	case pb.HandRanking_NaturalSixPairs:
		return pb.HandBonusType_BonusNaturalSixPairs
	case pb.HandRanking_NaturalThreeOfFlushes:
		return pb.HandBonusType_BonusNaturalThreeOfFlushes
	case pb.HandRanking_NaturalThreeStraights:
		return pb.HandBonusType_BonusNaturalThreeStraights
	}

	return pb.HandBonusType_None
}

func (h *Hand) CompareHand(h2 *Hand) *ComparisonResult {
	result := ComparisonResult{
		//WinType: entity.WinType_WIN_TYPE_UNSPECIFIED,
	}

	return &result
}

var kPc = "pc"

func CompareHand(ctx context.Context, h1, h2 *Hand) *ComparisonResult {
	// A (MB) vs
	//			B(MB) => case 1
	//			B(BL) => case 2
	//			B(BT) => case 3
	// A (BL) vs
	//			B(MB) => case 2
	//			B(BL) => case 4
	//			B(BT) => case 5
	// A (BT) vs
	//			B(MB) => case 3
	//			B(BL) => case 5
	//			B(BT) => case 6

	result := &ComparisonResult{}
	//count := ctx.Value(kPc).(int)
	h1.calculatePoint()
	h2.calculatePoint()

	// case 1
	if h1.IsNatural() && h2.IsNatural() {
		compareNaturalWithNatural(h1, h2, result)
		return result
	}

	// case 2
	if h1.IsNatural() && h2.IsMisSet() {
		compareNaturalWithMisset(h1, h2, result)
		return result
	}

	if h2.IsNatural() && h1.IsMisSet() {
		compareNaturalWithMisset(h2, h1, result)
		result.swap()
		return result
	}

	// case 3
	if h1.IsNatural() && h2.IsNormal() {
		compareNaturalWithNormal(h1, h2, result)
		return result
	}

	if h2.IsNatural() && h1.IsNormal() {
		compareNaturalWithNormal(h2, h1, result)
		result.swap()
		return result
	}

	// case 4
	if h1.IsMisSet() && h2.IsMisSet() {
		compareMissetWithMisset(h1, h2, result)
		return result
	}

	// case 5
	if h1.IsNormal() && h2.IsMisSet() {
		compareNormalWithMisset(h1, h2, result)
		return result
	}

	if h2.IsNormal() && h1.IsMisSet() {
		compareNormalWithMisset(h2, h1, result)
		result.swap()
		return result
	}

	// case 6
	if h1.IsNormal() && h2.IsNormal() {
		compareNormalWithNormal(h1, h2, result)
		return result
	}

	return result
}

//compareNaturalWithNatural
//case 1
func compareNaturalWithNatural(h1, h2 *Hand, result *ComparisonResult) {
	var score = 0
	if cmp := CompareHandPoint(h1.naturalPoint, h2.naturalPoint); cmp > 0 {
		score = mapNaturalPoint[rankingTypeToBonusType(h1.naturalPoint.rankingType)]
	} else if cmp < 0 {
		score = mapNaturalPoint[rankingTypeToBonusType(h2.naturalPoint.rankingType)]
	}

	result.r1.NaturalFactor = score
	result.r2.NaturalFactor = -score
}

//compareNaturalWithMisset
//case2
func compareNaturalWithMisset(h1, h2 *Hand, result *ComparisonResult) {
	score := mapNaturalPoint[rankingTypeToBonusType(h1.naturalPoint.rankingType)]
	result.r1.NaturalFactor = score
	result.r2.NaturalFactor = -score
}

//compareNaturalWithNormal
//case3
func compareNaturalWithNormal(h1, h2 *Hand, result *ComparisonResult) {
	score := mapNaturalPoint[rankingTypeToBonusType(h1.naturalPoint.rankingType)]
	result.r1.NaturalFactor = score
	result.r2.NaturalFactor = -score
}

//compareMissetWithMisset
//case4
func compareMissetWithMisset(h1, h2 *Hand, result *ComparisonResult) {
	// Don't need to do anything
}

//compareNormalWithMisset
func compareNormalWithMisset(h1, h2 *Hand, result *ComparisonResult) {
	bonusScoop := mapBonusPoint[pb.HandBonusType_Scoop]
	result.r1.ScoopFactor = bonusScoop
	result.r2.ScoopFactor = -bonusScoop

	// check special case bonus only
	if bonus, bonusScore := h1.frontHand.GetBonus(); bonus != pb.HandBonusType_None {
		result.r1.FrontFactor += bonusScore
		result.r2.FrontFactor += -bonusScore
	}

	if bonus, bonusScore := h1.middleHand.GetBonus(); bonus != pb.HandBonusType_None {
		result.r1.MiddleFactor += bonusScore
		result.r2.MiddleFactor += -bonusScore
	}

	if bonus, bonusScore := h1.backHand.GetBonus(); bonus != pb.HandBonusType_None {
		result.r1.BackFactor += bonusScore
		result.r2.BackFactor += -bonusScore
	}
}

//compareNormalWithNormal
func compareNormalWithNormal(h1, h2 *Hand, result *ComparisonResult) {
	// front hand
	if cmp := CompareHandPoint(h1.frontHand.Point, h2.frontHand.Point); cmp > 0 {
		if bonus, bonusScore := h1.frontHand.GetBonus(); bonus != pb.HandBonusType_None {
			result.r1.FrontBonusFactor = bonusScore
			result.r2.FrontBonusFactor = -bonusScore
		}

		result.r1.FrontFactor = 1
		result.r2.FrontFactor = -1
	} else if cmp < 0 {
		if bonus, bonusScore := h2.frontHand.GetBonus(); bonus != pb.HandBonusType_None {
			result.r2.FrontBonusFactor = bonusScore
			result.r1.FrontBonusFactor = -bonusScore
		}

		result.r2.FrontFactor = 1
		result.r1.FrontFactor = -1
	}

	// middle hand
	if cmp := CompareHandPoint(h1.middleHand.Point, h2.middleHand.Point); cmp > 0 {
		if bonus, bonusScore := h1.middleHand.GetBonus(); bonus != pb.HandBonusType_None {
			result.r1.MiddleBonusFactor = bonusScore
			result.r2.MiddleBonusFactor = -bonusScore
		}

		result.r1.MiddleFactor = 1
		result.r2.MiddleFactor = -1
	} else if cmp < 0 {
		if bonus, bonusScore := h2.middleHand.GetBonus(); bonus != pb.HandBonusType_None {
			result.r2.MiddleBonusFactor = bonusScore
			result.r1.MiddleBonusFactor = -bonusScore
		}

		result.r2.MiddleFactor = 1
		result.r1.MiddleFactor = -1
	}

	// backhand
	if cmp := CompareHandPoint(h1.backHand.Point, h2.backHand.Point); cmp > 0 {
		if bonus, bonusScore := h1.backHand.GetBonus(); bonus != pb.HandBonusType_None {
			result.r1.BackBonusFactor = bonusScore
			result.r2.BackBonusFactor = -bonusScore
		}

		result.r1.BackFactor = 1
		result.r2.BackFactor = -1
	} else if cmp < 0 {
		if bonus, bonusScore := h2.backHand.GetBonus(); bonus != pb.HandBonusType_None {
			result.r2.BackBonusFactor = bonusScore
			result.r1.MiddleBonusFactor = -bonusScore
		}

		result.r2.BackFactor = 1
		result.r1.BackFactor = -1
	}
}
