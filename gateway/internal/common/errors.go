package common

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	ForeignKeyViolation string = "23503"
	UniqueViolation     string = "23505"
)

const (
	NotValid         string = "not-valid"
	Unauthorized     string = "Unautorized"
	NotFound         string = "not-found"
	PermissionDenied string = "permission-denied"
	RequiredFields   string = "required-fields"
	TooManyReq       string = "too-many-requests"
	InternalError    string = "internal-server-error"
)

var (
	ErrRecordAlreadyExists  = errors.New("already exists")
	ErrEmailNotValid        = errors.New("email not valid")
	ErrPasswordNotValid     = errors.New("password not valid")
	ErrUserNotFound         = errors.New("user not found")
	ErrUserInactive         = errors.New("user is inactive")
	ErrUserActive           = errors.New("user is active")
	ErrUserNotVerified      = errors.New("user is not verified")
	ErrUserVerified         = errors.New("user has been verified")
	ErrTokenUsed            = errors.New("token has been used")
	ErrInvalidToken         = errors.New("token is invalid")
	ErrExpiredToken         = errors.New("token is expired")
	ErrInvalidOTP           = errors.New("otp is invalid")
	ErrSessionBlocked       = errors.New("blocked session")
	ErrSessionExpired       = errors.New("session is expired")
	ErrInternalError        = errors.New("internal server error")
	ErrSessionIncorrectUser = errors.New("incorrect user session")
)

func ConvertRPCError(err error) string {
	stats := status.Convert(err)

	message := stats.Message()
	return message
}

func ConvertRPCCodeError(err error, code codes.Code) bool {
	stats := status.Convert(err)
	return stats.Code() == code
}
