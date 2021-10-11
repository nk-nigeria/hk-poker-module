package main

import (
	"context"
	"database/sql"
	"github.com/ciaolink-game-platform/cgp-blackjack-module/api"
	"github.com/heroiclabs/nakama-common/runtime"
	"google.golang.org/protobuf/encoding/protojson"
)

func rpcGameList(marshaler *protojson.MarshalOptions, unmarshaler *protojson.UnmarshalOptions) func(context.Context, runtime.Logger, *sql.DB, runtime.NakamaModule, string) (string, error) {
	return func(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
		response, err := marshaler.Marshal(&api.GameListResponse{
			Games: []*api.Game{
				{
					Code:   "GAME1",
					Active: true,
				},
				{
					Code:   "GAME2",
					Active: true,
				},
				{
					Code:   "GAME3",
					Active: true,
				},
			},
		})
		if err != nil {
			logger.Error("error marshaling response payload: %v", err.Error())
			return "", errMarshal
		}

		return string(response), nil
	}
}
