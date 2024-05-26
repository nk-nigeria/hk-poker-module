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

	"github.com/ciaolink-game-platform/cgp-chinese-poker-module/api/presenter"
	"github.com/ciaolink-game-platform/cgp-chinese-poker-module/cgbdb"
	"github.com/ciaolink-game-platform/cgp-chinese-poker-module/entity"
	"github.com/ciaolink-game-platform/cgp-chinese-poker-module/pkg/packager"
	"github.com/ciaolink-game-platform/cgp-chinese-poker-module/usecase/engine"
	"github.com/ciaolink-game-platform/cgp-chinese-poker-module/usecase/processor"
	gsm "github.com/ciaolink-game-platform/cgp-chinese-poker-module/usecase/state_machine"
	pb "github.com/ciaolink-game-platform/cgp-common/proto"
	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/qmuntal/stateless"
	"google.golang.org/protobuf/encoding/protojson"
)

// Compile-time check to make sure all required functions are implemented.
var _ runtime.Match = &MatchHandler{}

type MatchHandler struct {
	processor processor.UseCase
	machine   gsm.UseCase
}

func (m *MatchHandler) MatchSignal(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, data string) (interface{}, string) {
	//panic("implement me")
	s := state.(*entity.MatchState)
	return s, ""
}

func NewMatchHandler(marshaler *protojson.MarshalOptions, unmarshaler *protojson.UnmarshalOptions) *MatchHandler {
	return &MatchHandler{
		processor: processor.NewMatchProcessor(marshaler, unmarshaler, engine.NewChinesePokerEngine()),
		machine:   gsm.NewGameStateMachine(),
	}
}

func (m *MatchHandler) GetState() stateless.State {
	return m.machine.MustState()
}

func (m *MatchHandler) MatchInit(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, params map[string]interface{}) (interface{}, int, string) {
	logger.Info("match init: %v", params)
	label, ok := params["data"].(string)
	if !ok {
		logger.WithField("params", params).Error("invalid match init parameter \"data\"")
		return nil, entity.TickRate, ""
	}
	matchInfo := &pb.Match{}
	err := entity.DefaulUnmarshaler.Unmarshal([]byte(label), matchInfo)
	if err != nil {
		logger.Error("match init json label failed ", err)
		return nil, entity.TickRate, ""
	}
	matchInfo.MatchId, _ = ctx.Value(runtime.RUNTIME_CTX_MATCH_ID).(string)
	labelJSON, err := entity.DefaultMarshaler.Marshal(matchInfo)

	if err != nil {
		logger.Error("match init json label failed ", err)
		return nil, entity.TickRate, ""
	}

	logger.Info("match init label= %s", string(labelJSON))

	matchState := entity.NewMathState(matchInfo)
	// init jp treasure
	jpTreasure, _ := cgbdb.GetJackpot(ctx, logger, db, entity.ModuleName)
	if jpTreasure != nil {
		matchState.SetJackpotTreasure(&pb.Jackpot{
			GameCode: jpTreasure.GetGameCode(),
			Chips:    jpTreasure.Chips,
		})
	}
	// fire idle event
	procPkg := packager.NewProcessorPackage(&matchState, m.processor, logger, nk, db, nil, nil, ctx)
	m.machine.TriggerIdle(packager.GetContextWithProcessorPackager(procPkg))

	return &matchState, entity.TickRate, string(labelJSON)
}

func (m *MatchHandler) MatchLoop(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, messages []runtime.MatchData) interface{} {
	s := state.(*entity.MatchState)

	err := m.machine.FireProcessEvent(packager.GetContextWithProcessorPackager(
		packager.NewProcessorPackage(
			s, m.processor,
			logger,
			nk,
			db,
			dispatcher,
			messages,
			ctx),
	))
	if err == presenter.ErrGameFinish {
		logger.Info("match need finish")

		return nil
	}

	return s
}

func (m *MatchHandler) MatchTerminate(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, graceSeconds int) interface{} {
	logger.Info("match terminate, state=%v")
	s := state.(*entity.MatchState)
	m.processor.ProcessMatchTerminate(ctx, logger, nk, db, dispatcher, s)
	return state
}
