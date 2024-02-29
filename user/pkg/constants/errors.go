package constants

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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
	ErrSessionIncorrectUser = errors.New("incorrect user session")
	ErrInternalError        = errors.New("internal server error")
	ErrRecordNotFound       = pgx.ErrNoRows
)

func ErrorCode(err error) string {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code
	}
	return ""
}
