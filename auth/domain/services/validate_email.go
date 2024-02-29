package services

import (
	"context"
	"errors"
	"time"

	"github.com/Foedie/foedie-server-v2/auth/domain/pb"
	"github.com/Foedie/foedie-server-v2/auth/internal/db"
	"github.com/Foedie/foedie-server-v2/auth/pkg/constants"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) ValidateEmail(ctx context.Context, req *pb.ValidateEmailParams) (*pb.VerifyEmailResponse, error) {
	token := req.GetToken()

	uid, err := tokenValid(token)

	if err != nil {
		return nil, status.Error(codes.Internal, constants.ErrInternalError.Error())
	}

	verifyEmail, errVerif := server.store.GetVerifyEmail(ctx, db.GetVerifyEmailParams{
		SecretKey: token,
		Uid:       uid,
	})

	if errors.Is(errVerif, constants.ErrRecordNotFound) {
		return nil, status.Error(codes.NotFound, constants.ErrUserNotFound.Error())
	}

	if verifyEmail.IsUsed {
		return nil, status.Error(codes.Aborted, constants.ErrTokenUsed.Error())
	}

	if verifyEmail.ExpiredAt.Time.Before(time.Now()) {
		return nil, status.Error(codes.Aborted, constants.ErrExpiredToken.Error())
	}

	err = server.store.ValidateEmailTx(ctx, uid, token)

	if err != nil {
		return nil, status.Error(codes.Internal, constants.ErrInternalError.Error())
	}

	return &pb.VerifyEmailResponse{
		Code:    "success",
		Message: "Email has been verified",
	}, nil
}
