package services

import (
	"context"
	"errors"

	"github.com/Foedie/foedie-server-v2/user/domain/pb"
	"github.com/Foedie/foedie-server-v2/user/internal/db"
	"github.com/Foedie/foedie-server-v2/user/pkg/constants"
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) UpdateUser(ctx context.Context, req *pb.UserParams) (*pb.SuccessResponse, error) {

	user, err := server.store.GetUser(ctx, req.GetUid())

	if errors.Is(err, constants.ErrRecordNotFound) {
		return nil, status.Error(codes.NotFound, constants.ErrUserNotFound.Error())
	}

	var userResp *pb.SuccessResponse
	userResp = &pb.SuccessResponse{
		Code:    "success",
		Message: "Successfully Updated!",
	}

	err = server.store.UpdateUserTx(ctx, user, db.UpdateUserParams{
		Uid:       req.GetUid(),
		Email:     req.GetEmail(),
		Username:  req.GetUsername(),
		FirstName: req.GetFirstName(),
		LastName: pgtype.Text{
			String: req.GetLastName(),
			Valid:  true,
		},
		PhoneNumber: req.GetPhoneNumber(),
	})

	if errors.Is(err, constants.ErrRecordNotFound) {
		return nil, status.Error(codes.NotFound, constants.ErrUserNotFound.Error())
	}

	if constants.ErrorCode(err) == constants.UniqueViolation {
		return nil, status.Error(codes.AlreadyExists, constants.ErrRecordAlreadyExists.Error())
	}

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if user.Email != req.GetEmail() {
		_, err := server.authSvc.VerifyEmail(req.GetEmail())
		if err != nil {
			code := status.Convert(err).Code()
			if code == codes.NotFound {
				return nil, status.Error(codes.NotFound, constants.ErrUserNotFound.Error())
			}
			if code == codes.Aborted {
				return nil, status.Error(codes.Aborted, constants.ErrUserVerified.Error())
			}
			return nil, status.Error(codes.Internal, err.Error())
		}

		userResp = &pb.SuccessResponse{
			Code:    "success",
			Message: "Email Verification link has been sent!",
		}
	}

	return userResp, nil
}
