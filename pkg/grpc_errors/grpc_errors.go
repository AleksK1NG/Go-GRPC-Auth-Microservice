package grpc_errors

import (
	"database/sql"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"strings"
)

// Parse error and get code
func ParseGRPCErrStatusCode(err error) codes.Code {
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return codes.NotFound
	case strings.Contains(err.Error(), "email") || strings.Contains(err.Error(), "password"):
		return codes.InvalidArgument
	}

	return codes.Internal
}
