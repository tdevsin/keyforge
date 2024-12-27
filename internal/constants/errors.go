package constants

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	StatusErrInvalidKey   = status.Errorf(codes.InvalidArgument, "Key is invalid")
	StatusErrInvalidValue = status.Errorf(codes.InvalidArgument, "Value is invalid")
	StatusErrKeyNotFound  = status.Errorf(codes.NotFound, "Key not found")
	StatusErrInternal     = status.Errorf(codes.Internal, "Some internal error occurred while processing your request")
)
