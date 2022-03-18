package api

import (
	"context"
	pb "github.com/ciaolink-game-platform/cgp-chinese-poker-module/proto"
	"github.com/heroiclabs/nakama-common/runtime"
)

func (m *MatchHandler) addChip(ctx context.Context, logger runtime.Logger, nk runtime.NakamaModule, userID string, amountChip int) {
	changeset := map[string]int64{
		"chips": int64(amountChip), // Add amountChip coins to the user's wallet.
	}
	metadata := map[string]interface{}{
		"game_topup": "topup",
	}

	_, _, err := nk.WalletUpdate(ctx, userID, changeset, metadata, true)
	if err != nil {
		logger.WithField("err", err).Error("Wallet update error.")
	}

}

func (m *MatchHandler) subtractChip(ctx context.Context, logger runtime.Logger, nk runtime.NakamaModule, userID string, amountChip int) {
	changeset := map[string]int64{
		"chips": -int64(amountChip), // Substract amountChip coins to the user's wallet.
	}
	metadata := map[string]interface{}{
		"game_topup": "topup",
	}

	_, _, err := nk.WalletUpdate(ctx, userID, changeset, metadata, true)
	if err != nil {
		logger.WithField("err", err).Error("Wallet update error.")
	}
}

func (m *MatchHandler) updateChipByResultGameFinish(ctx context.Context, logger runtime.Logger, nk runtime.NakamaModule, resultGame *pb.UpdateFinish) {
	walletUpdates := make([]*runtime.WalletUpdate, len(resultGame.Results))
	for _, result := range resultGame.Results {
		amountChip := int64(0)
		amountChip = 200*(result.FrontFactor+result.MiddleFactor+result.BackFactor) +
			(result.FrontBonus + result.MiddleBonus + result.BackBonus)
		changeset := map[string]int64{
			"chips": amountChip, // Substract amountChip coins to the user's wallet.
		}
		metadata := map[string]interface{}{
			"game_topup": "topup",
		}
		walletUpdates = append(walletUpdates, &runtime.WalletUpdate{
			UserID:    result.UserId,
			Changeset: changeset,
			Metadata:  metadata,
		})
	}

	_, err := nk.WalletsUpdate(ctx, walletUpdates, true)
	if err != nil {
		logger.WithField("err", err).Error("Wallets update error.")
	}
}
