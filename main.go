package main

import (
	"context"
	"database/sql"
	"github.com/ciaolink-game-platform/cgp-chinese-poker-module/entity"
	"github.com/heroiclabs/nakama-common/runtime"
	"google.golang.org/protobuf/encoding/protojson"
	"time"

	"github.com/ciaolink-game-platform/cgp-chinese-poker-module/api"
)

const (
	rpcIdGameList    = "list_game"
	rpcIdFindMatch   = "find_match"
	rpcIdCreateMatch = "create_match"
)

//noinspection GoUnusedExportedFunction
func InitModule(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, initializer runtime.Initializer) error {
	initStart := time.Now()

	marshaler := &protojson.MarshalOptions{
		UseEnumNumbers: true,
	}
	unmarshaler := &protojson.UnmarshalOptions{
		DiscardUnknown: false,
	}

	if err := initializer.RegisterRpc(rpcIdGameList, api.RpcGameList(marshaler, unmarshaler)); err != nil {
		return err
	}

	if err := initializer.RegisterRpc(rpcIdFindMatch, api.RpcFindMatch(marshaler, unmarshaler)); err != nil {
		return err
	}

	if err := initializer.RegisterRpc(rpcIdCreateMatch, api.RpcCreateMatch(marshaler, unmarshaler)); err != nil {
		return err
	}

	if err := initializer.RegisterMatch(entity.ModuleName, func(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule) (runtime.Match, error) {
		return api.NewMatchHandler(marshaler, unmarshaler), nil
	}); err != nil {
		return err
	}

	if err := api.RegisterSessionEvents(db, nk, initializer); err != nil {
		return err
	}

	logger.Info("Plugin loaded in '%d' msec.", time.Now().Sub(initStart).Milliseconds())
	return nil
}
