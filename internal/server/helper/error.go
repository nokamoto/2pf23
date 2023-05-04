package helper

import (
	"errors"

	"github.com/nokamoto/2pf23/internal/app"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ErrorOr returns a grpc error if err is not nil. Otherwise, returns res.
// If err is app.ErrInvalidArgument, returns codes.InvalidArgument.
// If err is app.ErrNotFound, returns codes.NotFound.
// Otherwise, returns codes.Unknown.
func ErrorOr[T any](logger *zap.Logger, res *T, err error) (*T, error) {
	if errors.Is(err, app.ErrInvalidArgument) {
		logger.Debug("invalid argument", zap.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if errors.Is(err, app.ErrNotFound) {
		logger.Debug("not found", zap.Error(err))
		return nil, status.Error(codes.NotFound, err.Error())
	}
	if err != nil {
		logger.Error("unknown error", zap.Error(err))
		return nil, status.Error(codes.Unknown, "unknown error")
	}
	return res, nil
}
