package grpc

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/typewriterco/p402/internal/problems"
)

// ToGRPCError converts a problems.Problem to a gRPC status error
func ToGRPCError(err error) error {
	if err == nil {
		return nil
	}

	var p problems.Problem
	if !errors.As(err, &p) {
		// Not a Problem error, return as Internal
		return status.Errorf(codes.Internal, "%v", err)
	}

	code := problemKindToGRPCCode(p.Kind)
	return status.Error(code, p.Detail)
}

// problemKindToGRPCCode maps problems.Kind to gRPC status codes
func problemKindToGRPCCode(kind problems.Kind) codes.Code {
	switch kind {
	case problems.Exist:
		return codes.AlreadyExists
	case problems.NotExist:
		return codes.NotFound
	case problems.Invalid, problems.Validation, problems.InvalidRequest:
		return codes.InvalidArgument
	case problems.Unauthenticated:
		return codes.Unauthenticated
	case problems.Unauthorized:
		return codes.PermissionDenied
	case problems.Private:
		return codes.PermissionDenied
	case problems.Database, problems.Internal, problems.Other, problems.IO, problems.Unanticipated:
		return codes.Internal
	case problems.BrokenLink:
		return codes.FailedPrecondition
	default:
		return codes.Internal
	}
}
