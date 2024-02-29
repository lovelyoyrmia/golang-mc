package services

import (
	"context"
	"errors"
	"time"

	"github.com/Foedie/foedie-server-v2/auth/domain/pb"
	"github.com/Foedie/foedie-server-v2/auth/internal/worker"
	"github.com/Foedie/foedie-server-v2/auth/pkg/constants"
	"github.com/hibiken/asynq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) VerifyEmail(ctx context.Context, req *pb.VerifyEmailParams) (*pb.VerifyEmailResponse, error) {
	user, err := server.store.GetUserByEmailOrUsername(ctx, req.GetEmail())

	if errors.Is(err, constants.ErrRecordNotFound) {
		return nil, status.Error(codes.NotFound, constants.ErrUserNotFound.Error())
	}

	if user.IsVerified.Bool {
		return nil, status.Error(codes.Aborted, constants.ErrUserVerified.Error())
	}

	taskPayload := &worker.PayloadSendVerifyEmail{
		Email: user.Email,
	}

	opts := []asynq.Option{
		asynq.MaxRetry(10),
		asynq.ProcessIn(10 * time.Second),
		asynq.Queue(worker.QueueCritical),
	}
	err = server.taskUser.UserTaskSendVerificationEmail(ctx, taskPayload, opts...)

	if err != nil {
		return nil, status.Error(codes.Internal, constants.ErrInternalError.Error())
	}

	return &pb.VerifyEmailResponse{
		Code:    "success",
		Message: "Email Verification link has been sent!",
	}, nil
}
