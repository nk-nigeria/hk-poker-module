// Copyright 2020 The Nakama Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/ciaolink-game-platform/cgp-chinese-poker-module/api/presenter"
	"github.com/ciaolink-game-platform/cgp-chinese-poker-module/entity"
	"github.com/ciaolink-game-platform/cgp-chinese-poker-module/pkg/packager"
	"github.com/ciaolink-game-platform/cgp-chinese-poker-module/usecase/chinese_poker"
	"github.com/ciaolink-game-platform/cgp-chinese-poker-module/usecase/game_state_machine"
	"github.com/ciaolink-game-platform/cgp-chinese-poker-module/usecase/processor"
	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/qmuntal/stateless"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
	tickRate = 5
)

// Compile-time check to make sure all required functions are implemented.
var _ runtime.Match = &MatchHandler{}

type MatchHandler struct {
	processor processor.UseCase
	machine   game_state_machine.UseCase
}

func NewMatchHandler(marshaler *protojson.MarshalOptions, unmarshaler *protojson.UnmarshalOptions) *MatchHandler {
	return &MatchHandler{
		processor: processor.NewMatchProcessor(marshaler, unmarshaler, chinese_poker.NewChinesePokerEngine()),
		machine:   game_state_machine.NewGameStateMachine(),
	}
}

func (m *MatchHandler) GetState() stateless.State {
	return m.machine.MustState()
}

func (m *MatchHandler) MatchInit(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, params map[string]interface{}) (interface{}, int, string) {
	logger.Info("match init: %v", params)
	bet, ok := params["bet"].(int32)
	if !ok {
		logger.Error("invalid match init parameter \"bet\"")
		return nil, 0, ""
	}

	name, ok := params["name"].(string)
	if !ok {
		logger.Error("invalid match init parameter \"name\"")
		return nil, 0, ""
	}

	password, ok := params["password"].(string)
	if !ok {
		logger.Error("invalid match init parameter \"password\"")
		return nil, 0, ""
	}

	label := &entity.MatchLabel{
		Open:     1,
		Bet:      bet,
		Code:     entity.ModuleName,
		Name:     name,
		Password: password,
		MaxSize:  entity.MaxPresences,
	}

	labelJSON, err := json.Marshal(label)
	if err != nil {
		logger.Error("match init json label failed ", err)
		return nil, tickRate, ""
	}

	logger.Info("match init label=", string(labelJSON))

	matchState := entity.NewMathState(label)

	// fire idle event
	procPkg := packager.NewProcessorPackage(&matchState, m.processor, logger, nil, nil)
	m.machine.TriggerIdle(packager.GetContextWithProcessorPackager(procPkg))

	return &matchState, tickRate, string(labelJSON)
}

func (m *MatchHandler) MatchLoop(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, messages []runtime.MatchData) interface{} {
	s := state.(*entity.MatchState)

	err := m.machine.FireProcessEvent(packager.GetContextWithProcessorPackager(packager.NewProcessorPackage(s, m.processor, logger, dispatcher, messages)))
	if err == presenter.ErrGameFinish {
		logger.Info("match need finish")
		return nil
	}

	return s
}

func (m *MatchHandler) MatchTerminate(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, graceSeconds int) interface{} {
	logger.Info("match terminate, state=%v")
	return state
}
