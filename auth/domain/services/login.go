package services

import (
	"context"
	"errors"
	"time"

	"github.com/Foedie/foedie-server-v2/auth/domain/pb"
	"github.com/Foedie/foedie-server-v2/auth/internal/db"
	"github.com/Foedie/foedie-server-v2/auth/pkg/constants"
	"github.com/Foedie/foedie-server-v2/auth/pkg/utils"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserParams) (*pb.LoginUserResponse, error) {
	mtd := extractMetadata(ctx)

	user, err := server.store.GetUserByEmailOrUsername(ctx, req.GetEmail())

	if errors.Is(err, constants.ErrRecordNotFound) {
		return nil, status.Error(codes.NotFound, constants.ErrUserNotFound.Error())
	}

	err = utils.ComparePassword(user.Password, req.GetPassword())

	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return nil, status.Error(codes.Unauthenticated, constants.ErrPasswordNotValid.Error())
	}

	if !user.IsVerified.Bool {
		return nil, status.Error(codes.Unauthenticated, constants.ErrUserNotVerified.Error())
	}

	if !user.IsActive.Bool {
		return nil, status.Error(codes.Unauthenticated, constants.ErrUserInactive.Error())
	}

	expiredAt := (time.Hour * 24) * 30
	refreshToken, payload, errRefreshToken := server.maker.GenerateToken(user.Uid, expiredAt)

	if errRefreshToken != nil {
		return nil, status.Error(codes.Internal, constants.ErrInternalError.Error())
	}

	accessToken, _, errAccessToken := server.maker.GenerateToken(user.Uid, time.Hour*24)

	if errAccessToken != nil {
		return nil, status.Error(codes.Internal, constants.ErrInternalError.Error())
	}

	refreshTokenParams := db.CreateRefreshTokenParams{
		ID:  payload.ID.String(),
		Uid: user.Uid,
		RefreshToken: pgtype.Text{
			String: refreshToken,
			Valid:  true,
		},
		UserAgent: pgtype.Text{
			String: mtd.UserAgent,
			Valid:  true,
		},
		ClientIp: pgtype.Text{
			String: mtd.ClientIP,
			Valid:  true,
		},
		ExpiredAt: pgtype.Timestamp{
			Time:  payload.ExpiredAt,
			Valid: true,
		},
	}

	session, errRefresh := server.store.CreateRefreshToken(ctx, refreshTokenParams)

	if errRefresh != nil {
		return nil, status.Error(codes.Internal, constants.ErrInternalError.Error())
	}

	err = server.store.UpdateUserLastLogin(ctx, user.Uid)

	if err != nil {
		return nil, status.Error(codes.Internal, constants.ErrInternalError.Error())
	}

	loginRes := &pb.LoginResponse{
		Token:        accessToken,
		RefreshToken: session.RefreshToken.String,
		Username:     user.Username,
		SecretKey:    user.SecretKey,
	}

	return &pb.LoginUserResponse{
		Code:    "success",
		Message: "Successfully Logged in !",
		Data:    loginRes,
	}, nil
}
