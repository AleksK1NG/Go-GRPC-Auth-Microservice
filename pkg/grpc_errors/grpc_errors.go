package grpc_errors

import (
	"context"
	"database/sql"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"strings"
)

var (
	ErrNotFound         = errors.New("Not found")
	ErrNoCtxMetaData    = errors.New("No ctx metadata")
	ErrInvalidSessionId = errors.New("Invalid session id")
)

// Parse error and get code
func ParseGRPCErrStatusCode(err error) codes.Code {
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return codes.NotFound
	case errors.Is(err, redis.Nil):
		return codes.NotFound
	case errors.Is(err, context.Canceled):
		return codes.Canceled
	case errors.Is(err, context.DeadlineExceeded):
		return codes.DeadlineExceeded
	case strings.Contains(err.Error(), "Validate"):
		return codes.InvalidArgument
	}
	return codes.Internal
}
