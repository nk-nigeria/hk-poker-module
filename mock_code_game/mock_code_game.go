package mockcodegame

import (
	"errors"
	"io"
	"net/http"

	pb "github.com/nakamaFramework/cgp-common/proto"
)

var MapMockCodeListCard = make(map[int][]*pb.Card)

func InitMapMockCodeListCard() {
	// 1. Sảnh đồng chất AKQJ 10,9,8,7,6,5,4,3,2
	// {
	// 	cards := []*pb.Card{
	// 		{Rank: pb.CardRank_RANK_2, Suit: pb.CardSuit_SUIT_SPADES}, {Rank: pb.CardRank_RANK_3, Suit: pb.CardSuit_SUIT_SPADES}, {Rank: pb.CardRank_RANK_4, Suit: pb.CardSuit_SUIT_SPADES},
	// 		{Rank: pb.CardRank_RANK_9, Suit: pb.CardSuit_SUIT_SPADES}, {Rank: pb.CardRank_RANK_8, Suit: pb.CardSuit_SUIT_SPADES}, {Rank: pb.CardRank_RANK_7, Suit: pb.CardSuit_SUIT_SPADES}, {Rank: pb.CardRank_RANK_6, Suit: pb.CardSuit_SUIT_SPADES}, {Rank: pb.CardRank_RANK_5, Suit: pb.CardSuit_SUIT_SPADES},
	// 		{Rank: pb.CardRank_RANK_10, Suit: pb.CardSuit_SUIT_SPADES}, {Rank: pb.CardRank_RANK_J, Suit: pb.CardSuit_SUIT_SPADES}, {Rank: pb.CardRank_RANK_Q, Suit: pb.CardSuit_SUIT_SPADES}, {Rank: pb.CardRank_RANK_K, Suit: pb.CardSuit_SUIT_SPADES}, {Rank: pb.CardRank_RANK_A, Suit: pb.CardSuit_SUIT_SPADES},
	// 	}
	// 	MapMockCodeListCard[1] = cards
	// }
	// // 2. Sảnh rồng AKQJ rô -10 tépm 98765 cơ 432 tép
	// {
	// 	cards := []*pb.Card{
	// 		{Rank: pb.CardRank_RANK_2, Suit: pb.CardSuit_SUIT_CLUBS}, {Rank: pb.CardRank_RANK_3, Suit: pb.CardSuit_SUIT_CLUBS}, {Rank: pb.CardRank_RANK_4, Suit: pb.CardSuit_SUIT_CLUBS},
	// 		{Rank: pb.CardRank_RANK_9, Suit: pb.CardSuit_SUIT_HEARTS}, {Rank: pb.CardRank_RANK_8, Suit: pb.CardSuit_SUIT_HEARTS}, {Rank: pb.CardRank_RANK_7, Suit: pb.CardSuit_SUIT_HEARTS}, {Rank: pb.CardRank_RANK_6, Suit: pb.CardSuit_SUIT_HEARTS}, {Rank: pb.CardRank_RANK_5, Suit: pb.CardSuit_SUIT_HEARTS},
	// 		{Rank: pb.CardRank_RANK_10, Suit: pb.CardSuit_SUIT_CLUBS}, {Rank: pb.CardRank_RANK_J, Suit: pb.CardSuit_SUIT_DIAMONDS}, {Rank: pb.CardRank_RANK_Q, Suit: pb.CardSuit_SUIT_DIAMONDS}, {Rank: pb.CardRank_RANK_K, Suit: pb.CardSuit_SUIT_DIAMONDS}, {Rank: pb.CardRank_RANK_A, Suit: pb.CardSuit_SUIT_DIAMONDS},
	// 	}
	// 	MapMockCodeListCard[2] = cards
	// }
	// // 3. 6 đôi:  5 rô 6 tép 6 rô 7 tép 7 rô, 3 cơ 3 rô 4 cơ 4 rô 5 tép, 2 cơ 2 rô, 8 rô
	// {
	// 	cards := []*pb.Card{
	// 		{Rank: pb.CardRank_RANK_2, Suit: pb.CardSuit_SUIT_HEARTS}, {Rank: pb.CardRank_RANK_2, Suit: pb.CardSuit_SUIT_DIAMONDS}, {Rank: pb.CardRank_RANK_8, Suit: pb.CardSuit_SUIT_DIAMONDS},
	// 		{Rank: pb.CardRank_RANK_3, Suit: pb.CardSuit_SUIT_HEARTS}, {Rank: pb.CardRank_RANK_3, Suit: pb.CardSuit_SUIT_DIAMONDS}, {Rank: pb.CardRank_RANK_4, Suit: pb.CardSuit_SUIT_HEARTS}, {Rank: pb.CardRank_RANK_4, Suit: pb.CardSuit_SUIT_DIAMONDS}, {Rank: pb.CardRank_RANK_5, Suit: pb.CardSuit_SUIT_CLUBS},
	// 		{Rank: pb.CardRank_RANK_7, Suit: pb.CardSuit_SUIT_HEARTS}, {Rank: pb.CardRank_RANK_7, Suit: pb.CardSuit_SUIT_DIAMONDS}, {Rank: pb.CardRank_RANK_6, Suit: pb.CardSuit_SUIT_HEARTS}, {Rank: pb.CardRank_RANK_6, Suit: pb.CardSuit_SUIT_DIAMONDS}, {Rank: pb.CardRank_RANK_5, Suit: pb.CardSuit_SUIT_DIAMONDS},
	// 	}
	// 	MapMockCodeListCard[3] = cards
	// }
	// // 4. 3 sảnh AKQJ10 tépm 87654 bích 4 rô 3 tép 2 tép
	// {
	// 	cards := []*pb.Card{
	// 		{Rank: pb.CardRank_RANK_4, Suit: pb.CardSuit_SUIT_CLUBS}, {Rank: pb.CardRank_RANK_3, Suit: pb.CardSuit_SUIT_CLUBS}, {Rank: pb.CardRank_RANK_2, Suit: pb.CardSuit_SUIT_CLUBS},
	// 		{Rank: pb.CardRank_RANK_4, Suit: pb.CardSuit_SUIT_SPADES}, {Rank: pb.CardRank_RANK_8, Suit: pb.CardSuit_SUIT_SPADES}, {Rank: pb.CardRank_RANK_7, Suit: pb.CardSuit_SUIT_SPADES}, {Rank: pb.CardRank_RANK_6, Suit: pb.CardSuit_SUIT_SPADES}, {Rank: pb.CardRank_RANK_5, Suit: pb.CardSuit_SUIT_SPADES},
	// 		{Rank: pb.CardRank_RANK_10, Suit: pb.CardSuit_SUIT_CLUBS}, {Rank: pb.CardRank_RANK_J, Suit: pb.CardSuit_SUIT_CLUBS}, {Rank: pb.CardRank_RANK_Q, Suit: pb.CardSuit_SUIT_CLUBS}, {Rank: pb.CardRank_RANK_K, Suit: pb.CardSuit_SUIT_CLUBS}, {Rank: pb.CardRank_RANK_A, Suit: pb.CardSuit_SUIT_CLUBS},
	// 	}
	// 	MapMockCodeListCard[4] = cards
	// }
	// // 5. Same color 2,3,5,6,8 cơ 78910 rô QKJ cơ
	// {
	// 	cards := []*pb.Card{
	// 		{Rank: pb.CardRank_RANK_J, Suit: pb.CardSuit_SUIT_HEARTS}, {Rank: pb.CardRank_RANK_Q, Suit: pb.CardSuit_SUIT_HEARTS}, {Rank: pb.CardRank_RANK_K, Suit: pb.CardSuit_SUIT_HEARTS},
	// 		{Rank: pb.CardRank_RANK_7, Suit: pb.CardSuit_SUIT_DIAMONDS}, {Rank: pb.CardRank_RANK_8, Suit: pb.CardSuit_SUIT_DIAMONDS}, {Rank: pb.CardRank_RANK_9, Suit: pb.CardSuit_SUIT_DIAMONDS}, {Rank: pb.CardRank_RANK_10, Suit: pb.CardSuit_SUIT_DIAMONDS}, {Rank: pb.CardRank_RANK_6, Suit: pb.CardSuit_SUIT_DIAMONDS},
	// 		{Rank: pb.CardRank_RANK_2, Suit: pb.CardSuit_SUIT_HEARTS}, {Rank: pb.CardRank_RANK_3, Suit: pb.CardSuit_SUIT_HEARTS}, {Rank: pb.CardRank_RANK_5, Suit: pb.CardSuit_SUIT_HEARTS}, {Rank: pb.CardRank_RANK_6, Suit: pb.CardSuit_SUIT_HEARTS}, {Rank: pb.CardRank_RANK_8, Suit: pb.CardSuit_SUIT_HEARTS},
	// 	}
	// 	MapMockCodeListCard[5] = cards
	// }
	// // 6. 3 bộ đồng chất A,4,6,8,10 tép, 3,5,7,9,j rô JQK tép
	// {
	// 	cards := []*pb.Card{
	// 		{Rank: pb.CardRank_RANK_J, Suit: pb.CardSuit_SUIT_CLUBS}, {Rank: pb.CardRank_RANK_Q, Suit: pb.CardSuit_SUIT_CLUBS}, {Rank: pb.CardRank_RANK_K, Suit: pb.CardSuit_SUIT_CLUBS},
	// 		{Rank: pb.CardRank_RANK_3, Suit: pb.CardSuit_SUIT_DIAMONDS}, {Rank: pb.CardRank_RANK_5, Suit: pb.CardSuit_SUIT_DIAMONDS}, {Rank: pb.CardRank_RANK_7, Suit: pb.CardSuit_SUIT_DIAMONDS}, {Rank: pb.CardRank_RANK_9, Suit: pb.CardSuit_SUIT_DIAMONDS}, {Rank: pb.CardRank_RANK_J, Suit: pb.CardSuit_SUIT_DIAMONDS},
	// 		{Rank: pb.CardRank_RANK_A, Suit: pb.CardSuit_SUIT_CLUBS}, {Rank: pb.CardRank_RANK_4, Suit: pb.CardSuit_SUIT_CLUBS}, {Rank: pb.CardRank_RANK_6, Suit: pb.CardSuit_SUIT_CLUBS}, {Rank: pb.CardRank_RANK_8, Suit: pb.CardSuit_SUIT_CLUBS}, {Rank: pb.CardRank_RANK_10, Suit: pb.CardSuit_SUIT_CLUBS},
	// 	}
	// 	MapMockCodeListCard[6] = cards
	// }
	// // 7. Thùng phá sảnh 5,4,3,2,A bích | 5,4,3,2 cơ A tép, 7,8,9 bích
	// {
	// 	cards := []*pb.Card{
	// 		{Rank: pb.CardRank_RANK_7, Suit: pb.CardSuit_SUIT_SPADES}, {Rank: pb.CardRank_RANK_8, Suit: pb.CardSuit_SUIT_SPADES}, {Rank: pb.CardRank_RANK_9, Suit: pb.CardSuit_SUIT_SPADES},
	// 		{Rank: pb.CardRank_RANK_5, Suit: pb.CardSuit_SUIT_HEARTS}, {Rank: pb.CardRank_RANK_4, Suit: pb.CardSuit_SUIT_HEARTS}, {Rank: pb.CardRank_RANK_3, Suit: pb.CardSuit_SUIT_HEARTS}, {Rank: pb.CardRank_RANK_2, Suit: pb.CardSuit_SUIT_HEARTS}, {Rank: pb.CardRank_RANK_A, Suit: pb.CardSuit_SUIT_CLUBS},
	// 		{Rank: pb.CardRank_RANK_5, Suit: pb.CardSuit_SUIT_SPADES}, {Rank: pb.CardRank_RANK_4, Suit: pb.CardSuit_SUIT_SPADES}, {Rank: pb.CardRank_RANK_3, Suit: pb.CardSuit_SUIT_SPADES}, {Rank: pb.CardRank_RANK_2, Suit: pb.CardSuit_SUIT_SPADES}, {Rank: pb.CardRank_RANK_A, Suit: pb.CardSuit_SUIT_SPADES},
	// 	}
	// 	MapMockCodeListCard[7] = cards
	// }
	// AAAA2 88855 444
	// {
	// 	cards := []*pb.Card{
	// 		{Rank: pb.CardRank_RANK_4, Suit: pb.CardSuit_SUIT_SPADES}, {Rank: pb.CardRank_RANK_4, Suit: pb.CardSuit_SUIT_HEARTS}, {Rank: pb.CardRank_RANK_4, Suit: pb.CardSuit_SUIT_DIAMONDS},
	// 		{Rank: pb.CardRank_RANK_8, Suit: pb.CardSuit_SUIT_SPADES}, {Rank: pb.CardRank_RANK_8, Suit: pb.CardSuit_SUIT_HEARTS}, {Rank: pb.CardRank_RANK_8, Suit: pb.CardSuit_SUIT_DIAMONDS}, {Rank: pb.CardRank_RANK_5, Suit: pb.CardSuit_SUIT_HEARTS}, {Rank: pb.CardRank_RANK_5, Suit: pb.CardSuit_SUIT_CLUBS},
	// 		{Rank: pb.CardRank_RANK_A, Suit: pb.CardSuit_SUIT_SPADES}, {Rank: pb.CardRank_RANK_A, Suit: pb.CardSuit_SUIT_DIAMONDS}, {Rank: pb.CardRank_RANK_A, Suit: pb.CardSuit_SUIT_CLUBS}, {Rank: pb.CardRank_RANK_A, Suit: pb.CardSuit_SUIT_HEARTS}, {Rank: pb.CardRank_RANK_2, Suit: pb.CardSuit_SUIT_SPADES},
	// 	}
	// 	MapMockCodeListCard[1] = cards
	// }
	// // AkQJ10 bich 88887 384
	// {
	// 	cards := []*pb.Card{
	// 		{Rank: pb.CardRank_RANK_4, Suit: pb.CardSuit_SUIT_SPADES}, {Rank: pb.CardRank_RANK_3, Suit: pb.CardSuit_SUIT_HEARTS}, {Rank: pb.CardRank_RANK_8, Suit: pb.CardSuit_SUIT_DIAMONDS},
	// 		{Rank: pb.CardRank_RANK_8, Suit: pb.CardSuit_SUIT_SPADES}, {Rank: pb.CardRank_RANK_8, Suit: pb.CardSuit_SUIT_HEARTS}, {Rank: pb.CardRank_RANK_8, Suit: pb.CardSuit_SUIT_DIAMONDS}, {Rank: pb.CardRank_RANK_8, Suit: pb.CardSuit_SUIT_CLUBS}, {Rank: pb.CardRank_RANK_7, Suit: pb.CardSuit_SUIT_CLUBS},
	// 		{Rank: pb.CardRank_RANK_10, Suit: pb.CardSuit_SUIT_SPADES}, {Rank: pb.CardRank_RANK_J, Suit: pb.CardSuit_SUIT_SPADES}, {Rank: pb.CardRank_RANK_Q, Suit: pb.CardSuit_SUIT_SPADES}, {Rank: pb.CardRank_RANK_K, Suit: pb.CardSuit_SUIT_SPADES}, {Rank: pb.CardRank_RANK_A, Suit: pb.CardSuit_SUIT_SPADES},
	// 	}
	// 	MapMockCodeListCard[2] = cards
	// }
	// // AkQJ10 cơ 76543 nhép 694
	// {
	// 	cards := []*pb.Card{
	// 		{Rank: pb.CardRank_RANK_6, Suit: pb.CardSuit_SUIT_HEARTS}, {Rank: pb.CardRank_RANK_9, Suit: pb.CardSuit_SUIT_HEARTS}, {Rank: pb.CardRank_RANK_4, Suit: pb.CardSuit_SUIT_DIAMONDS},
	// 		{Rank: pb.CardRank_RANK_7, Suit: pb.CardSuit_SUIT_SPADES}, {Rank: pb.CardRank_RANK_6, Suit: pb.CardSuit_SUIT_SPADES}, {Rank: pb.CardRank_RANK_5, Suit: pb.CardSuit_SUIT_SPADES}, {Rank: pb.CardRank_RANK_4, Suit: pb.CardSuit_SUIT_SPADES}, {Rank: pb.CardRank_RANK_3, Suit: pb.CardSuit_SUIT_SPADES},
	// 		{Rank: pb.CardRank_RANK_10, Suit: pb.CardSuit_SUIT_HEARTS}, {Rank: pb.CardRank_RANK_J, Suit: pb.CardSuit_SUIT_HEARTS}, {Rank: pb.CardRank_RANK_Q, Suit: pb.CardSuit_SUIT_HEARTS}, {Rank: pb.CardRank_RANK_K, Suit: pb.CardSuit_SUIT_HEARTS}, {Rank: pb.CardRank_RANK_A, Suit: pb.CardSuit_SUIT_HEARTS},
	// 	}
	// 	MapMockCodeListCard[3] = cards
	// }
	// jp
	{
		cards := []*pb.Card{
			{
				Rank: pb.CardRank_RANK_A, Suit: pb.CardSuit_SUIT_SPADES}, {Rank: pb.CardRank_RANK_A, Suit: pb.CardSuit_SUIT_CLUBS}, {Rank: pb.CardRank_RANK_J, Suit: pb.CardSuit_SUIT_SPADES},
			{Rank: pb.CardRank_RANK_10, Suit: pb.CardSuit_SUIT_SPADES}, {Rank: pb.CardRank_RANK_Q, Suit: pb.CardSuit_SUIT_SPADES}, {Rank: pb.CardRank_RANK_K, Suit: pb.CardSuit_SUIT_SPADES}, {Rank: pb.CardRank_RANK_2, Suit: pb.CardSuit_SUIT_SPADES}, {Rank: pb.CardRank_RANK_3, Suit: pb.CardSuit_SUIT_SPADES},
			{Rank: pb.CardRank_RANK_4, Suit: pb.CardSuit_SUIT_SPADES}, {Rank: pb.CardRank_RANK_5, Suit: pb.CardSuit_SUIT_SPADES}, {Rank: pb.CardRank_RANK_6, Suit: pb.CardSuit_SUIT_SPADES}, {Rank: pb.CardRank_RANK_7, Suit: pb.CardSuit_SUIT_SPADES}, {Rank: pb.CardRank_RANK_8, Suit: pb.CardSuit_SUIT_SPADES},
		}
		MapMockCodeListCard[1] = cards
	}
}

func downloadFile(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode > 300 {
		return nil, errors.New("status not ok")
	}
	if resp.Body == nil {
		return nil, errors.New("body is nil")
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}
