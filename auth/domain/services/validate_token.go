package services

import (
	"context"
	"errors"
	"net/http"

	"github.com/Foedie/foedie-server-v2/auth/domain/pb"
	"github.com/Foedie/foedie-server-v2/auth/pkg/constants"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) ValidateToken(ctx context.Context, req *pb.ValidateTokenParams) (*pb.ValidateTokenResponse, error) {
	payload, err := server.maker.VerifyToken(req.GetToken())

	if errors.Is(err, constants.ErrInvalidToken) {
		return nil, status.Error(codes.Unauthenticated, constants.ErrInvalidToken.Error())
	}

	if errors.Is(err, constants.ErrExpiredToken) {
		return nil, status.Error(codes.Unauthenticated, constants.ErrExpiredToken.Error())
	}

	user, err := server.store.GetUserByUid(ctx, payload.UID)

	if errors.Is(err, constants.ErrRecordNotFound) {
		return nil, status.Error(codes.NotFound, constants.ErrUserNotFound.Error())
	}

	if !user.IsActive.Bool {
		return nil, status.Error(codes.Unauthenticated, constants.ErrUserInactive.Error())
	}

	if !user.IsVerified.Bool {
		return nil, status.Error(codes.Unauthenticated, constants.ErrUserNotVerified.Error())
	}

	if err != nil {
		return nil, status.Error(codes.Internal, constants.ErrInternalError.Error())
	}

	return &pb.ValidateTokenResponse{
		Error:  "",
		Status: http.StatusOK,
		Uid:    user.Uid,
	}, nil
}
