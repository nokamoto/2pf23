package helper

import (
	"github.com/bufbuild/connect-go"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
)

type connectMatcher[T any] struct {
	x *connect.Request[T]
	f func(*connect.Request[T]) proto.Message
}

func (m *connectMatcher[T]) Matches(x interface{}) bool {
	y, ok := x.(*connect.Request[T])
	if !ok {
		return false
	}
	return proto.Equal(m.f(m.x), m.f(y))
}

func (m *connectMatcher[T]) String() string {
	return prototext.Format(m.f(m.x))
}

func ConnectEqual[T any](x *T) *connectMatcher[T] {
	return &connectMatcher[T]{
		x: &connect.Request[T]{Msg: x},
		f: func(x *connect.Request[T]) proto.Message {
			if x == nil || x.Msg == nil {
				return nil
			}
			var m interface{} = x.Msg
			return m.(proto.Message)
		},
	}
}
