package mock

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
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
	Name   string               `json:"name"`
	Input  InputChinsePokerMock `json:"input"`
	Output pb.UpdateFinish      `json:"output"`
}

func TestChinsePokerAllMock(t *testing.T) {
	path := "./chinese_poker_mock"
	files, err := ioutil.ReadDir(path)
	if err != nil {
		t.Fatal(err)
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}
		nameFile := f.Name()
		if !strings.HasSuffix(nameFile, ".json") {
			continue
		}
		// t.Logf(f.Name())
		fullFilePath := filepath.Join(path, f.Name())
		RunTestChinsePokerMock(fullFilePath, t)
	}
}

func RunTestChinsePokerMock(fileMock string, t *testing.T) {
	data, err := os.ReadFile(fileMock) // just pass the file name
	if err != nil {
		t.Fatalf("Error read file mock %s , err %s", fileMock, err.Error())
	}
	cpMock := &ChinsePokerMock{}
	err = json.Unmarshal(data, &cpMock)
	if err != nil {
		t.Fatalf("Parse file mock %s err %s", fileMock, err.Error())
	}
	t.Logf("Run test for mock %s", cpMock.Name)

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
	sort.Slice(cpMock.Output.Bonuses, func(i, j int) bool {
		a := cpMock.Output.Bonuses[i]
		b := cpMock.Output.Bonuses[j]
		if a.Win < b.Win {
			return true
		}
		if a.Lose < b.Lose {
			return true
		}
		return a.Type.Number() < b.Type.Number()
	})
	sort.Slice(result.Bonuses, func(i, j int) bool {
		a := result.Bonuses[i]
		b := result.Bonuses[j]
		if a.Win < b.Win {
			return true
		}
		if a.Lose < b.Lose {
			return true
		}
		return a.Type.Number() < b.Type.Number()
	})
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
			fmt.Sprintf("%s - %s", cpMock.Name, "type point -1(misset) 0(normal) 1(natual) "))

		assert.Equal(t, expectResult.ScoreResult.FrontBonusFactor,
			actualResult.ScoreResult.FrontBonusFactor,
			fmt.Sprintf("%s - %s", cpMock.Name, "front bonus factor"))
		assert.Equal(t, expectResult.ScoreResult.MiddleBonusFactor,
			actualResult.ScoreResult.MiddleBonusFactor,
			fmt.Sprintf("%s - %s", cpMock.Name, "mid bonus factor"))
		assert.Equal(t, expectResult.ScoreResult.BackBonusFactor,
			actualResult.ScoreResult.BackBonusFactor,
			fmt.Sprintf("%s - %s", cpMock.Name, "back bonus factor"))

		assert.Equal(t, expectResult.ScoreResult.FrontFactor,
			actualResult.ScoreResult.FrontFactor,
			fmt.Sprintf("%s - %s", cpMock.Name, "front factor"))
		assert.Equal(t, expectResult.ScoreResult.MiddleFactor,
			actualResult.ScoreResult.MiddleFactor,
			fmt.Sprintf("%s - %s", cpMock.Name, "mid factor"))
		assert.Equal(t, expectResult.ScoreResult.BackFactor,
			actualResult.ScoreResult.BackFactor,
			fmt.Sprintf("%s - %s", cpMock.Name, "back factor"))

		assert.Equal(t, expectResult.ScoreResult.NaturalFactor,
			actualResult.ScoreResult.NaturalFactor,
			fmt.Sprintf("%s - %s", cpMock.Name, "NaturalFactor"))
		assert.Equal(t, expectResult.ScoreResult.NumHandWin,
			actualResult.ScoreResult.NumHandWin,
			fmt.Sprintf("%s - %s", cpMock.Name, "NumHandWin"))
		// assert.Equal(t, expectResult.ScoreResult.Scoop,
		// 	actualResult.ScoreResult.Scoop,
		// 	"Scoop")

	}

	assert.Equal(t, len(cpMock.Output.Bonuses), len(result.Bonuses), "len arr bonus")
	//sort bonus by user id

	for idx, expect := range cpMock.Output.Bonuses {
		actual := result.Bonuses[idx]
		assert.Equal(t, expect.Win, actual.Win, fmt.Sprintf("%s - %s", cpMock.Name, "user id win"))
		assert.Equal(t, expect.Lose, actual.Lose, fmt.Sprintf("%s - %s", cpMock.Name, "user id lose"))
		assert.Equal(t, expect.Factor, actual.Factor, fmt.Sprintf("%s - %s", cpMock.Name, "factor"))
		assert.Equal(t, expect.Type, actual.Type, fmt.Sprintf("%s - %s", cpMock.Name, "Type"))
	}
	t.Logf("%v", result)
}
func TestChinsePokerMock(t *testing.T) {
	fileMock := "./chinese_poker_mock/normal-win-all.json"
	RunTestChinsePokerMock(fileMock, t)
}
