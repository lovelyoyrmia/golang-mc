package services

import (
	"context"
	"errors"
	"time"

	"github.com/Foedie/foedie-server-v2/auth/domain/pb"
	"github.com/Foedie/foedie-server-v2/auth/internal/db"
	"github.com/Foedie/foedie-server-v2/auth/internal/worker"
	"github.com/Foedie/foedie-server-v2/auth/pkg/constants"
	"github.com/hibiken/asynq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) RecoverAccount(ctx context.Context, req *pb.RecoverAccountParams) (*pb.RecoverAccountResponse, error) {
	user, err := server.store.GetUserByEmailAndUsername(ctx, db.GetUserByEmailAndUsernameParams{
		Username: req.GetUsername(),
		Email:    req.GetEmail(),
	})

	if errors.Is(err, constants.ErrRecordNotFound) {
		return nil, status.Error(codes.NotFound, constants.ErrUserNotFound.Error())
	}

	if user.IsActive.Bool {
		return nil, status.Error(codes.Aborted, constants.ErrUserActive.Error())
	}

	if !user.IsVerified.Bool {
		return nil, status.Error(codes.Aborted, constants.ErrUserNotVerified.Error())
	}

	taskPayload := &worker.PayloadSendRecoverAccount{
		Email:    user.Email,
		Username: user.Username,
	}

	opts := []asynq.Option{
		asynq.MaxRetry(10),
		asynq.ProcessIn(10 * time.Second),
		asynq.Queue(worker.QueueCritical),
	}

	err = server.taskUser.UserTaskSendRecoverAccount(ctx, taskPayload, opts...)

	if err != nil {
		return nil, status.Error(codes.Internal, constants.ErrInternalError.Error())
	}

	return &pb.RecoverAccountResponse{
		Code:    "success",
		Message: "Please check your email to recover your account",
	}, nil

}
