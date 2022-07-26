package hand

import (
	"reflect"
	"testing"

	"github.com/ciaolink-game-platform/cgp-chinese-poker-module/entity"
	pb "github.com/ciaolink-game-platform/cgp-chinese-poker-module/proto"
)

func mockHandNatural1() (*Hand, error) {
	return NewHandFromPb(&pb.ListCard{
		Cards: []*pb.Card{
			{
				Rank: pb.CardRank_RANK_6,
				Suit: pb.CardSuit_SUIT_SPADES,
			},
			{
				Rank: pb.CardRank_RANK_A,
				Suit: pb.CardSuit_SUIT_SPADES,
			},
			{
				Rank: pb.CardRank_RANK_8,
				Suit: pb.CardSuit_SUIT_SPADES,
			},

			{
				Rank: pb.CardRank_RANK_9,
				Suit: pb.CardSuit_SUIT_HEARTS,
			},
			{
				Rank: pb.CardRank_RANK_7,
				Suit: pb.CardSuit_SUIT_HEARTS,
			},
			{
				Rank: pb.CardRank_RANK_7,
				Suit: pb.CardSuit_SUIT_DIAMONDS,
			},
			{
				Rank: pb.CardRank_RANK_10,
				Suit: pb.CardSuit_SUIT_HEARTS,
			},
			{
				Rank: pb.CardRank_RANK_4,
				Suit: pb.CardSuit_SUIT_HEARTS,
			},

			{
				Rank: pb.CardRank_RANK_5,
				Suit: pb.CardSuit_SUIT_HEARTS,
			},
			{
				Rank: pb.CardRank_RANK_5,
				Suit: pb.CardSuit_SUIT_CLUBS,
			},
			{
				Rank: pb.CardRank_RANK_3,
				Suit: pb.CardSuit_SUIT_HEARTS,
			},
			{
				Rank: pb.CardRank_RANK_3,
				Suit: pb.CardSuit_SUIT_CLUBS,
			},
			{
				Rank: pb.CardRank_RANK_3,
				Suit: pb.CardSuit_SUIT_DIAMONDS,
			},
		},
	})
}

func mockHandNatural2() (*Hand, error) {
	return NewHandFromPb(&pb.ListCard{
		Cards: []*pb.Card{
			{
				Rank: pb.CardRank_RANK_9,
				Suit: pb.CardSuit_SUIT_SPADES,
			},
			{
				Rank: pb.CardRank_RANK_K,
				Suit: pb.CardSuit_SUIT_SPADES,
			},
			{
				Rank: pb.CardRank_RANK_4,
				Suit: pb.CardSuit_SUIT_SPADES,
			},

			{
				Rank: pb.CardRank_RANK_8,
				Suit: pb.CardSuit_SUIT_SPADES,
			},
			{
				Rank: pb.CardRank_RANK_4,
				Suit: pb.CardSuit_SUIT_DIAMONDS,
			},
			{
				Rank: pb.CardRank_RANK_J,
				Suit: pb.CardSuit_SUIT_DIAMONDS,
			},
			{
				Rank: pb.CardRank_RANK_5,
				Suit: pb.CardSuit_SUIT_SPADES,
			},
			{
				Rank: pb.CardRank_RANK_10,
				Suit: pb.CardSuit_SUIT_CLUBS,
			},
			{
				Rank: pb.CardRank_RANK_J,
				Suit: pb.CardSuit_SUIT_DIAMONDS,
			},

			{
				Rank: pb.CardRank_RANK_A,
				Suit: pb.CardSuit_SUIT_DIAMONDS,
			},
			{
				Rank: pb.CardRank_RANK_10,
				Suit: pb.CardSuit_SUIT_SPADES,
			},
			{
				Rank: pb.CardRank_RANK_Q,
				Suit: pb.CardSuit_SUIT_SPADES,
			},
			{
				Rank: pb.CardRank_RANK_K,
				Suit: pb.CardSuit_SUIT_HEARTS,
			},
		},
	})
}

func TestCompareHand(t *testing.T) {

	h1, _ := mockHandNatural1()
	h2, _ := mockHandNatural2()
	ctx := NewCompareContext(2)
	result := CompareHand(ctx, h1, h2)
	r1 := pb.ComparisonResult{
		ScoreResult: &pb.ScoreResult{},
		PointResult: &pb.PointResult{},
	}
	r2 := pb.ComparisonResult{
		ScoreResult: &pb.ScoreResult{},
		PointResult: &pb.PointResult{},
	}

	ProcessCompareResult(ctx, &r1, result.GetR1())
	ProcessCompareResult(ctx, &r2, result.GetR2())
	t.Logf("result %v", result)
	t.Logf("result %v", result.bonuses)
	t.Logf("r1 %v", r1.ScoreResult)
	t.Logf("r1 %v", r1.PointResult)
	t.Logf("r2 %v", r2.ScoreResult)
	t.Logf("r2 %v", r2.PointResult)

}

func TestHand_CompareHand(t *testing.T) {
	type fields struct {
		cards        entity.ListCard
		ranking      pb.HandRanking
		frontHand    *ChildHand
		middleHand   *ChildHand
		backHand     *ChildHand
		naturalPoint *HandPoint
		pointType    pb.PointType
		calculated   bool
		owner        string
	}
	type args struct {
		h2 *Hand
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *ComparisonResult
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Hand{
				cards:        tt.fields.cards,
				ranking:      tt.fields.ranking,
				frontHand:    tt.fields.frontHand,
				middleHand:   tt.fields.middleHand,
				backHand:     tt.fields.backHand,
				naturalPoint: tt.fields.naturalPoint,
				pointType:    tt.fields.pointType,
				calculated:   tt.fields.calculated,
				owner:        tt.fields.owner,
			}
			if got := h.CompareHand(tt.args.h2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Hand.CompareHand() = %v, want %v", got, tt.want)
			}
		})
	}
}
