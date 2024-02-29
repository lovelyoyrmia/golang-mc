package services

import (
	"context"
	"errors"

	"github.com/Foedie/foedie-server-v2/user/domain/pb"
	"github.com/Foedie/foedie-server-v2/user/pkg/constants"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) DeleteUser(ctx context.Context, req *pb.UserUidParams) (*pb.SuccessResponse, error) {

	err := server.store.DeleteUser(ctx, req.GetUid())

	if errors.Is(err, constants.ErrRecordNotFound) {
		return nil, status.Error(codes.NotFound, constants.ErrUserNotFound.Error())
	}

	if err != nil {
		return nil, status.Error(codes.Internal, constants.ErrInternalError.Error())
	}

	return &pb.SuccessResponse{
		Code:    "success",
		Message: "Success",
	}, nil
}
