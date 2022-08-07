package mock

import (
	"encoding/json"
	"os"
	"sort"
	"testing"

	"github.com/ciaolink-game-platform/cgp-chinese-poker-module/entity"
	pb "github.com/ciaolink-game-platform/cgp-chinese-poker-module/proto"
	"github.com/ciaolink-game-platform/cgp-chinese-poker-module/usecase/engine"
	"github.com/emirpasic/gods/maps/linkedhashmap"
	"github.com/stretchr/testify/assert"
)

type MockCard struct {
	UserId string `json:"userId"`
	Card   string `json:"card"`
}

type InputChinsePokerMock struct {
	Cards []MockCard `json:"cards"`
}

type ChinsePokerMock struct {
	Input  InputChinsePokerMock `json:"input"`
	Output pb.UpdateFinish      `json:"output"`
}

func TestChinsePokerMock(t *testing.T) {
	fileMock := "/home/sondq/Documents/myspace/cgb-chinese-poker-module/mock/chinese_poker_mock/chinese_poker_mock_1.json"
	data, err := os.ReadFile(fileMock) // just pass the file name
	if err != nil {
		t.Fatalf("Error read file mock %s , err %s", fileMock, err.Error())
	}
	cpMock := &ChinsePokerMock{}
	err = json.Unmarshal(data, &cpMock)
	if err != nil {
		t.Fatalf("Parse file mock %s err %s", fileMock, err.Error())
	}

	processor := engine.NewChinesePokerEngine()
	presense := linkedhashmap.New()
	for _, u := range cpMock.Input.Cards {
		presense.Put(u.UserId, nil)
	}
	state := &entity.MatchState{
		Presences:        presense,
		PlayingPresences: presense,
		OrganizeCards:    make(map[string]*pb.ListCard),
		Cards:            make(map[string]*pb.ListCard),
	}
	for _, u := range cpMock.Input.Cards {
		listCard := &pb.ListCard{
			Cards: ParseListCard(u.Card),
		}
		processor.Organize(state, u.UserId, listCard)
	}
	result := processor.Finish(state)
	mapExpectResult := make(map[string]*pb.ComparisonResult)
	for _, r := range cpMock.Output.Results {
		mapExpectResult[r.GetUserId()] = r
	}
	mapActualResult := make(map[string]*pb.ComparisonResult)
	for _, r := range result.Results {
		mapActualResult[r.GetUserId()] = r
	}
	for _, u := range cpMock.Input.Cards {
		userId := u.UserId
		expectResult := mapExpectResult[userId]
		actualResult := mapActualResult[userId]
		assert.Equal(t, expectResult.PointResult.Type,
			actualResult.PointResult.Type,
			"type point -1(misset) 0(normal) 1(natual), ")

		assert.Equal(t, expectResult.ScoreResult.FrontBonusFactor,
			actualResult.ScoreResult.FrontBonusFactor,
			"front bonus factor")
		assert.Equal(t, expectResult.ScoreResult.MiddleBonusFactor,
			actualResult.ScoreResult.MiddleBonusFactor,
			"mid bonus factor")
		assert.Equal(t, expectResult.ScoreResult.BackBonusFactor,
			actualResult.ScoreResult.BackBonusFactor,
			"back bonus factor")

		assert.Equal(t, expectResult.ScoreResult.FrontBonusFactor,
			actualResult.ScoreResult.FrontFactor,
			"front factor")
		assert.Equal(t, expectResult.ScoreResult.MiddleFactor,
			actualResult.ScoreResult.MiddleFactor,
			"mid factor")
		assert.Equal(t, expectResult.ScoreResult.BackFactor,
			actualResult.ScoreResult.BackFactor,
			"back factor")

		assert.Equal(t, expectResult.ScoreResult.NaturalFactor,
			actualResult.ScoreResult.NaturalFactor,
			"NaturalFactor")
		assert.Equal(t, expectResult.ScoreResult.NumHandWin,
			actualResult.ScoreResult.NumHandWin,
			"NumHandWin")
		// assert.Equal(t, expectResult.ScoreResult.Scoop,
		// 	actualResult.ScoreResult.Scoop,
		// 	"Scoop")

	}

	assert.Equal(t, len(cpMock.Output.Bonuses), len(result.Bonuses), "len arr bonus")
	//sort bonus by user id
	sort.Slice(cpMock.Output.Bonuses, func(i, j int) bool {
		a := cpMock.Output.Bonuses[i]
		b := cpMock.Output.Bonuses[j]
		if a.Win < b.Win {
			return true
		}
		if a.Win > b.Win {
			return false
		}
		return a.Lose < b.Lose
	})
	sort.Slice(result.Bonuses, func(i, j int) bool {
		a := result.Bonuses[i]
		b := result.Bonuses[j]
		if a.Win < b.Win {
			return true
		}
		if a.Win > b.Win {
			return false
		}
		return a.Lose < b.Lose
	})

	for idx, expect := range cpMock.Output.Bonuses {
		actual := result.Bonuses[idx]
		assert.Equal(t, expect.Win, actual.Win, "user id win")
		assert.Equal(t, expect.Lose, actual.Lose, "user id lose")
		assert.Equal(t, expect.Factor, actual.Factor, "factor")
		assert.Equal(t, expect.Type, actual.Type, "Type")
	}
	t.Logf("%v", result)
}
