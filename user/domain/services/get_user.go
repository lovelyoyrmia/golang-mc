package services

import (
	"context"
	"errors"

	"github.com/Foedie/foedie-server-v2/user/domain/pb"
	"github.com/Foedie/foedie-server-v2/user/pkg/constants"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) GetUser(ctx context.Context, req *pb.UserUidParams) (*pb.UserResponse, error) {

	user, err := server.store.GetUser(ctx, req.GetUid())

	if errors.Is(err, constants.ErrRecordNotFound) {
		return nil, status.Error(codes.NotFound, constants.ErrUserNotFound.Error())
	}

	if err != nil {
		return nil, status.Error(codes.Internal, constants.ErrInternalError.Error())
	}

	return &pb.UserResponse{
		Code:    "success",
		Message: "Success",
		Data: &pb.User{
			Uid:         user.Uid,
			Email:       user.Email,
			Username:    user.Username,
			FirstName:   user.FirstName,
			PhoneNumber: user.PhoneNumber,
			LastName:    user.LastName.String,
			IsVerified:  user.IsVerified.Bool,
			LastLogin:   user.LastLogin.Time.String(),
			OtpVerified: user.OtpVerified.Bool,
			OtpEnabled:  user.OtpEnabled.Bool,
			OtpAuthUrl:  user.OtpUrl.String,
		},
	}, nil
}
