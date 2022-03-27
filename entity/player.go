package entity

import (
	"context"
	"strconv"

	pb "github.com/ciaolink-game-platform/cgp-chinese-poker-module/proto"
	"github.com/heroiclabs/nakama-common/runtime"
)

type ArrPbPlayer []*pb.Player

func NewPlayer(presence runtime.Presence) *pb.Player {
	p := pb.Player{
		Id:       presence.GetUserId(),
		UserName: presence.GetUsername(),
	}
	return &p
}

func NewListPlayer(presences []runtime.Presence) ArrPbPlayer {
	listPlayer := make([]*pb.Player, 0, len(presences))
	for _, presense := range presences {
		p := NewPlayer(presense)
		listPlayer = append(listPlayer, p)
	}
	return listPlayer
}

func (pb ArrPbPlayer) ReadWallet(ctx context.Context, nk runtime.NakamaModule, logger runtime.Logger) error {
	listUserId := make([]string, 0, len(pb))
	for _, player := range pb {
		listUserId = append(listUserId, player.Id)
	}
	wallets, err := ReadWalletUsers(ctx, nk, logger, listUserId...)
	if err != nil {
		return err
	}
	mapWallet := make(map[string]Wallet)
	for _, w := range wallets {
		mapWallet[w.UserId] = w
	}
	for i, player := range pb {
		player.Wallet = strconv.FormatInt(mapWallet[player.Id].Chips, 10)
		pb[i] = player
	}
	return nil
}
