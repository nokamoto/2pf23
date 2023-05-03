package helper

import (
	"github.com/golang/mock/gomock"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
)

type protoMatcher struct {
	x proto.Message
}

func (m *protoMatcher) Matches(x interface{}) bool {
	y, ok := x.(proto.Message)
	if !ok {
		return false
	}
	return proto.Equal(m.x, y)
}

func (m *protoMatcher) String() string {
	return prototext.Format(m.x)
}

// ProtoEqual returns a gomock.Matcher for proto.Message.
func ProtoEqual(x proto.Message) gomock.Matcher {
	return &protoMatcher{x}
}
