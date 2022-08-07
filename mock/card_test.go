package mock

import (
	"testing"

	pb "github.com/ciaolink-game-platform/cgp-chinese-poker-module/proto"
	"google.golang.org/protobuf/encoding/protojson"
)

func TestParseMockCard(t *testing.T) {
	fileMock := "/home/sondq/Documents/myspace/cgb-chinese-poker-module/mock/mock_card/natural_special.txt"
	list := ParseMockCard(fileMock)
	x := pb.UpdateFinish{}
	pp := &protojson.MarshalOptions{
		EmitUnpopulated: true,
		UseEnumNumbers:  true,
	}
	data, _ := pp.Marshal(&x)
	t.Logf("%s", string(data))
	t.Logf("%v", list)

}
