package grpc

import (
	"errors"
	"github.com/Royal17x/flagr/backend/internal/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func domainErrToGRPC(err error) error {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return status.Error(codes.NotFound, "not found")
	case errors.Is(err, domain.ErrAlreadyExists):
		return status.Error(codes.AlreadyExists, "already exists")
	case errors.Is(err, domain.ErrInvalidInput):
		return status.Error(codes.InvalidArgument, "invalid input")
	case errors.Is(err, domain.ErrUnauthorized):
		return status.Error(codes.Unauthenticated, "unauthorized")
	default:
		return status.Error(codes.Internal, "internal error")
	}
}
