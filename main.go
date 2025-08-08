package main

import (
	"context"
	"database/sql"
	"time"

	"github.com/nk-nigeria/cgp-common/bot"
	"github.com/nk-nigeria/cgp-common/define"
	"github.com/nk-nigeria/hk-poker-module/constant"
	"github.com/nk-nigeria/hk-poker-module/message_queue"
	mockcodegame "github.com/nk-nigeria/hk-poker-module/mock_code_game"

	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/nk-nigeria/hk-poker-module/api"
	"github.com/nk-nigeria/hk-poker-module/entity"
	_ "golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
	rpcIdGameList    = "list_game"
	rpcIdFindMatch   = "find_match"
	rpcIdCreateMatch = "create_match"
)

// noinspection GoUnusedExportedFunction
func InitModule(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, initializer runtime.Initializer) error {
	initStart := time.Now()
	bot.LoadBotsInfo(ctx, nk, db)
	message_queue.InitNatsService(logger, constant.NastEndpoint, &protojson.MarshalOptions{})
	mockcodegame.InitMapMockCodeListCard()
	if err := initializer.RegisterMatch(entity.ModuleName, func(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule) (runtime.Match, error) {
		return api.NewMatchHandler(entity.DefaultMarshaler, entity.DefaulUnmarshaler), nil
	}); err != nil {
		return err
	}

	if err := api.RegisterSessionEvents(db, nk, initializer); err != nil {
		return err
	}
	entity.BotLoader = bot.NewBotLoader(db, define.ChinesePoker.String(), 100000000)

	logger.Info("Plugin loaded in '%d' msec.", time.Now().Sub(initStart).Milliseconds())
	return nil
}
