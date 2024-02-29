package services

import (
	"context"
	"time"

	"github.com/Foedie/foedie-server-v2/auth/domain/pb"
	"github.com/Foedie/foedie-server-v2/auth/internal/db"
	"github.com/Foedie/foedie-server-v2/auth/internal/worker"
	"github.com/Foedie/foedie-server-v2/auth/pkg/constants"
	"github.com/Foedie/foedie-server-v2/auth/pkg/utils"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func toCreateUserParams(req *pb.CreateUserParams) (*db.CreateUserParams, error) {
	var uid string

	if req.GetUid() == "" {
		uid = uuid.NewString()
	} else {
		uid = req.GetUid()
	}

	if !utils.IsEmailValid(req.GetEmail()) {
		return nil, constants.ErrEmailNotValid
	}

	if req.GetPassword() != req.GetConfirmPassword() {
		return nil, constants.ErrPasswordNotValid
	}

	password, err := utils.HashPassword(req.GetPassword())

	if err != nil {
		return nil, err
	}

	secretKey, err := utils.Encrypt(uid)

	if err != nil {
		return nil, err
	}

	return &db.CreateUserParams{
		Uid:         uid,
		PhoneNumber: req.GetPhoneNumber(),
		Email:       req.GetEmail(),
		Username:    req.GetUsername(),
		FirstName:   req.GetFirstName(),
		LastName: pgtype.Text{
			String: req.GetLastName(),
			Valid:  true,
		},
		Password:  password,
		SecretKey: secretKey,
		IsActive: pgtype.Bool{
			Bool:  true,
			Valid: true,
		},
		IsVerified: pgtype.Bool{
			Bool:  false,
			Valid: true,
		},
	}, nil
}

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserParams) (*pb.CreateUserResponse, error) {

	userParams, err := toCreateUserParams(req)

	if err != nil {
		return nil, status.Error(codes.Internal, constants.ErrInternalError.Error())
	}

	arg := db.UserTxParams{
		UserParams: *userParams,
		AfterCreate: func(user db.User) error {
			taskPayload := &worker.PayloadSendVerifyEmail{
				Email: user.Email,
			}

			opts := []asynq.Option{
				asynq.MaxRetry(10),
				asynq.ProcessIn(10 * time.Second),
				asynq.Queue(worker.QueueCritical),
			}
			return server.taskUser.UserTaskSendVerificationEmail(ctx, taskPayload, opts...)
		},
	}

	err = server.store.CreateUserTx(ctx, arg)

	if constants.ErrorCode(err) == constants.UniqueViolation {
		return nil, status.Error(codes.AlreadyExists, constants.ErrRecordAlreadyExists.Error())
	}

	if err != nil {
		return nil, status.Error(codes.Internal, constants.ErrInternalError.Error())
	}

	return &pb.CreateUserResponse{
		Code:    "success",
		Message: "Successfully Created",
	}, nil
}
