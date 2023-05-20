package helper

import (
	"errors"
	"testing"

	"github.com/bufbuild/connect-go"
	"github.com/nokamoto/2pf23/internal/app"
	"github.com/nokamoto/2pf23/internal/util/helper"
	"go.uber.org/zap/zaptest"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestErrorOr(t *testing.T) {
	res := connect.NewResponse(&emptypb.Empty{})

	tests := []struct {
		name string
		err  error
		res  *connect.Response[emptypb.Empty]
		code connect.Code
	}{
		{
			name: "ok",
			err:  nil,
			res:  res,
			code: 0,
		},
		{
			name: "invalid argument",
			err:  app.ErrInvalidArgument,
			code: connect.CodeInvalidArgument,
		},
		{
			name: "not found",
			err:  app.ErrNotFound,
			code: connect.CodeNotFound,
		},
		{
			name: "unknown error",
			err:  errors.New("unknown error"),
			code: connect.CodeUnknown,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := zaptest.NewLogger(t)
			got, err := ErrorOr(logger, res, tt.err)
			if got != tt.res {
				t.Errorf("ErrorOr() got = %v, want %v", got, tt.res)
			}
			if code := helper.CodeOf(err); code != tt.code {
				t.Errorf("ErrorOr() code = %v, want %v", code, tt.code)
			}
		})
	}
}
