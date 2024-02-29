package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/Foedie/foedie-server-v2/gateway/data"
	"github.com/Foedie/foedie-server-v2/gateway/internal/common"
	"github.com/Foedie/foedie-server-v2/gateway/internal/models"
	"github.com/Foedie/foedie-server-v2/gateway/pkg/auth/pb"
)

type AuthMiddlewareConfig struct {
	svc *ServiceClient
}

func InitAuthMiddleware(svc *ServiceClient) AuthMiddlewareConfig {
	return AuthMiddlewareConfig{svc}
}

func (c *AuthMiddlewareConfig) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tokenString := checkAuthorization(w, r)

		if tokenString == "" {
			return
		}

		res, err := c.svc.Client.ValidateToken(context.Background(), &pb.ValidateTokenParams{
			Token: tokenString,
		})

		if err == nil {
			ctx := context.WithValue(r.Context(), models.KeyUser{}, res.Uid)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		if errMsg := common.ConvertRPCError(err); errMsg == common.ErrExpiredToken.Error() {
			errRes := &data.FailureResponse{
				Code:   common.Unauthorized,
				Detail: "Token is expired",
				Status: false,
			}
			data.Response(w, errRes, http.StatusUnauthorized)
			return
		}

		if errMsg := common.ConvertRPCError(err); errMsg == common.ErrInvalidToken.Error() {
			errRes := &data.FailureResponse{
				Code:   common.Unauthorized,
				Detail: "Token is invalid",
				Status: false,
			}
			data.Response(w, errRes, http.StatusUnauthorized)
			return
		}

		if errMsg := common.ConvertRPCError(err); errMsg == common.ErrUserNotFound.Error() {
			errRes := &data.FailureResponse{
				Code:   common.Unauthorized,
				Detail: "User not found",
				Status: false,
			}
			data.Response(w, errRes, http.StatusUnauthorized)
			return
		}

		if errMsg := common.ConvertRPCError(err); errMsg == common.ErrUserInactive.Error() {
			errRes := &data.FailureResponse{
				Code:   common.Unauthorized,
				Detail: "User is inactive",
				Status: false,
			}
			data.Response(w, errRes, http.StatusUnauthorized)
			return
		}

		if errMsg := common.ConvertRPCError(err); errMsg == common.ErrUserNotVerified.Error() {
			errRes := &data.FailureResponse{
				Code:   common.Unauthorized,
				Detail: "User is not verified",
				Status: false,
			}
			data.Response(w, errRes, http.StatusUnauthorized)
			return
		}

		errRes := &data.FailureResponse{
			Code:   common.InternalError,
			Detail: err.Error(),
			Status: false,
		}
		data.Response(w, errRes, http.StatusInternalServerError)
	})
}

func checkAuthorization(w http.ResponseWriter, r *http.Request) string {
	authHeader := r.Header.Get("Authorization")

	if authHeader == "" {
		errRes := &data.FailureResponse{
			Code:   common.Unauthorized,
			Detail: "Authorization was required",
			Status: false,
		}
		data.Response(w, errRes, http.StatusUnauthorized)
		return ""
	}

	if !strings.HasPrefix(authHeader, "Bearer") {
		errRes := &data.FailureResponse{
			Code:   common.Unauthorized,
			Detail: "Bearer token was required",
			Status: false,
		}
		data.Response(w, errRes, http.StatusUnauthorized)
		return ""
	}

	splitHeader := strings.Split(authHeader, " ")
	if len(splitHeader) < 2 {
		errRes := &data.FailureResponse{
			Code:   common.Unauthorized,
			Detail: "Bearer token was required",
			Status: false,
		}
		data.Response(w, errRes, http.StatusUnauthorized)
		return ""
	}

	tokenString := splitHeader[1]
	return tokenString
}
