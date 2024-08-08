package engine

import (
	"github.com/nakamaFramework/cgp-chinese-poker-module/entity"
	pb "github.com/nakamaFramework/cgp-common/proto"
)

type UseCase interface {
	NewGame(s *entity.MatchState) error
	Deal(s *entity.MatchState) error
	Organize(s *entity.MatchState, presence string, cards *pb.ListCard) error
	Combine(s *entity.MatchState, presence string) error
	Finish(s *entity.MatchState) *pb.UpdateFinish

	AISortCard(cards []*pb.Card) []*pb.Card
}
