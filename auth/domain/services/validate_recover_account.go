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

func (server *Server) ValidateRecoverAccount(ctx context.Context, req *pb.ValidateRecoverAccountParams) (*pb.RecoverAccountResponse, error) {
	token := req.GetToken()

	uid, err := tokenValid(token)

	if err != nil {
		return nil, status.Error(codes.Internal, constants.ErrInternalError.Error())
	}

	recoverAccount, err := server.store.GetRecoverAccount(ctx, db.GetRecoverAccountParams{
		SecretKey: token,
		Uid:       uid,
	})

	if errors.Is(err, constants.ErrRecordNotFound) {
		return nil, status.Error(codes.NotFound, constants.ErrUserNotFound.Error())
	}

	if recoverAccount.IsUsed {
		return nil, status.Error(codes.Aborted, constants.ErrTokenUsed.Error())
	}

	if recoverAccount.ExpiredAt.Time.Before(time.Now()) {
		return nil, status.Error(codes.Aborted, constants.ErrExpiredToken.Error())
	}

	err = server.store.ValidateRecoverAccountTx(ctx, uid, token)

	if err != nil {
		return nil, status.Error(codes.Internal, constants.ErrInternalError.Error())
	}

	return &pb.RecoverAccountResponse{
		Code:    "success",
		Message: "Account Recovered!",
	}, nil
}
