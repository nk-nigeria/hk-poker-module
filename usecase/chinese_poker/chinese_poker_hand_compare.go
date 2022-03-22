package chinese_poker

import "context"

func calculateCompareBonus(result *ComparisonResult) *ComparisonResult {
	//t := result.WinType

	//if t&entity.WIN_TYPE_WIN_FRONT_THREE_OF_A_KIND != 0 {
	//	result.FrontBonusFactor += int64(entity.GetWinFactorBonus(entity.WIN_TYPE_WIN_FRONT_THREE_OF_A_KIND))
	//}
	//if t&entity.WIN_TYPE_WIN_MID_FULL_HOUSE != 0 {
	//	result.MiddleBonusFactor += int64(entity.GetWinFactorBonus(entity.WIN_TYPE_WIN_MID_FULL_HOUSE))
	//}
	//if t&entity.WIN_TYPE_WIN_BACK_FOUR_OF_A_KIND != 0 {
	//	result.BackBonusFactor += int64(entity.GetWinFactorBonus(entity.WIN_TYPE_WIN_BACK_FOUR_OF_A_KIND))
	//}
	//if t&entity.WIN_TYPE_WIN_MID_FOUR_OF_A_KIND != 0 {
	//	result.MiddleBonusFactor += int64(entity.GetWinFactorBonus(entity.WIN_TYPE_WIN_MID_FOUR_OF_A_KIND))
	//}
	//if t&entity.WIN_TYPE_WIN_BACK_STRAIGHT_FLUSH != 0 {
	//	result.BackBonusFactor += int64(entity.GetWinFactorBonus(entity.WIN_TYPE_WIN_BACK_STRAIGHT_FLUSH))
	//}
	//if t&entity.WIN_TYPE_WIN_MID_STRAIGHT_FLUSH != 0 {
	//	result.MiddleBonusFactor += int64(entity.GetWinFactorBonus(entity.WIN_TYPE_WIN_MID_STRAIGHT_FLUSH))
	//}
	//
	//// Scoop check
	//if result.FrontFactor > 0 && result.MiddleFactor > 0 && result.BackFactor > 0 {
	//	result.ScoopFactor = 1
	//
	//	result.ScoopBonus = int64(entity.GetWinFactorBonus(entity.WIN_TYPE_WIN_SCOOP))
	//}

	return result
}

func (h *Hand) CompareHand(h2 *Hand) *ComparisonResult {
	result := ComparisonResult{
		//WinType: entity.WinType_WIN_TYPE_UNSPECIFIED,
	}

	return &result
}

var kPc = "pc"

func CompareHand(ctx context.Context, h1, h2 *Hand, result *ComparisonResult) {
	// A (MB) vs
	//			B(MB) => case 1
	//			B(BL) => case
	//			B(BT) =>
	// A (BL) vs
	//			B(MB) =>
	//			B(BL) =>
	//			B(BT) =>
	// A (BT) vs
	//			B(MB) =>
	//			B(BL) =>
	//			B(BT) =>

	//count := ctx.Value(kPc).(int)
	h1.calculatePoint()
	h2.calculatePoint()

	// case 1
	if h1.IsNatural() && h2.IsNatural() {
		compareNaturalWithNatural(h1, h2, result)
	}

	// case 2
	if h1.IsNatural() && h2.IsMisSet() {
		compareNaturalWithMisset(h1, h2, result)
	}

	if h2.IsNatural() && h1.IsMisSet() {
		compareNaturalWithMisset(h2, h1, result)
	}
}

func compareNaturalWithNatural(h1, h2 *Hand, result *ComparisonResult) *ComparisonResult {

	return nil
}

func compareNaturalWithMisset(h1, h2 *Hand, result *ComparisonResult) *ComparisonResult {
	return nil
}

func compareNaturalWithNormal(h1, h2 *Hand, result *ComparisonResult) *ComparisonResult {
	return nil
}

func compareMissetWithMisset(h1, h2 *Hand, result *ComparisonResult) *ComparisonResult {
	return nil
}

func compareMissetWithNormal(h1, h2 *Hand, result *ComparisonResult) *ComparisonResult {
	return nil
}

func compareNormalWithNormal(h1, h2 *Hand, result *ComparisonResult) *ComparisonResult {
	//result := h1.CompareHand(h2)
	//if result.WinType != 0 {
	//	return result
	//}

	////  chi dau
	//result.BackFactor = int64(h1.backHand.CompareHand(h2.backHand))
	//if result.BackFactor > 0 {
	//	r := h1.backHand.Point.rankingType
	//	switch r {
	//	case pb.HandRanking_FourOfAKind:
	//		result.WinType |= entity.WIN_TYPE_WIN_BACK_FOUR_OF_A_KIND
	//	case pb.HandRanking_StraightFlush:
	//		result.WinType |= entity.WIN_TYPE_WIN_BACK_STRAIGHT_FLUSH
	//	}
	//
	//}
	//// chi giua
	//result.MiddleFactor = int64(h1.middleHand.CompareHand(h2.middleHand))
	//if result.MiddleFactor > 0 {
	//	r := h1.middleHand.Point.rankingType
	//	switch r {
	//	case pb.HandRanking_FullHouse:
	//		result.WinType |= entity.WIN_TYPE_WIN_MID_FULL_HOUSE
	//	case pb.HandRanking_FourOfAKind:
	//		if h1.backHand.Point.rankingType == pb.HandRanking_FourOfAKind {
	//			result.WinType |= entity.WIN_TYPE_WIN_MID_FOUR_OF_A_KIND
	//		}
	//	case pb.HandRanking_StraightFlush:
	//		if h1.backHand.Point.rankingType == pb.HandRanking_StraightFlush {
	//			result.WinType |= entity.WIN_TYPE_WIN_MID_STRAIGHT_FLUSH
	//		}
	//	}
	//}
	//// chi cuoi
	//result.FrontFactor = int64(h1.frontHand.CompareHand(h2.frontHand))
	//if result.FrontBonusFactor > 0 && h1.frontHand.Point.rankingType == pb.HandRanking_ThreeOfAKind {
	//	result.WinType |= entity.WIN_TYPE_WIN_FRONT_THREE_OF_A_KIND
	//}

	calculateCompareBonus(result)

	return result
}
