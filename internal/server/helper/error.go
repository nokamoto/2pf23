package helper

import (
	"errors"

	"github.com/bufbuild/connect-go"
	"github.com/nokamoto/2pf23/internal/app"
	"go.uber.org/zap"
)

// ErrorOr returns a grpc error if err is not nil. Otherwise, returns res.
// If err is app.ErrInvalidArgument, returns codes.InvalidArgument.
// If err is app.ErrNotFound, returns codes.NotFound.
// Otherwise, returns codes.Unknown.
func ErrorOr[T any](logger *zap.Logger, res *connect.Response[T], err error) (*connect.Response[T], error) {
	if errors.Is(err, app.ErrInvalidArgument) {
		logger.Debug("invalid argument", zap.Error(err))
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}
	if errors.Is(err, app.ErrNotFound) {
		logger.Debug("not found", zap.Error(err))
		return nil, connect.NewError(connect.CodeNotFound, err)
	}
	if err != nil {
		logger.Error("unknown error", zap.Error(err))
		return nil, connect.NewError(connect.CodeUnknown, errors.New("unknown error"))
	}
	return res, nil
}

// CodeOf returns connect.CodeOf(err) if err is not nil. Otherwise, returns 0 for OK.
// connect.CodeOf(err) returns codes.Unknown if err is nil, but it is not useful for testing.
func CodeOf(err error) connect.Code {
	if err == nil {
		return 0
	}
	return connect.CodeOf(err)
}

// GetMsg returns v.Msg if res is not nil. Otherwise, returns nil.
//
// Original grpc-go code implements nil check in generated code to access message field and able to get default value,
// but connect-go does not because [connect.Request].Msg panics if it is nil.
//
// ```
// var foo *T = nil
// foo.GetCluster() // ok
//
// var bar connect.Request[T] = nil
// bar.Msg.GetCluster() // panic
// ```
//
// To avoid this panic, this function is introduced.
//
// ```
// var bar connect.Request[T] = nil
// GetMsg(bar).GetCluster() // ok
// ```
func GetMsg[T any](v *connect.Request[T]) *T {
	if v == nil {
		return nil
	}
	return v.Msg
}

// GetResponseMsg returns v.Msg if res is not nil. Otherwise, returns nil.
// See GetMsg for details.
func GetResponseMsg[T any](v *connect.Response[T]) *T {
	if v == nil {
		return nil
	}
	return v.Msg
}
