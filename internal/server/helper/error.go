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
