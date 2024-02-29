package routes

import (
	"net/http"

	"github.com/Foedie/foedie-server-v2/gateway/data"
	"github.com/Foedie/foedie-server-v2/gateway/internal/common"
	"github.com/Foedie/foedie-server-v2/gateway/pkg/auth/pb"
	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc/codes"
)

func (route *AuthRoutes) ValidateEmail(rw http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")

	if !isTokenValid(token) {
		errRes := data.FailureResponse{
			Detail: "Invalid token",
			Code:   common.NotFound,
			Status: false,
		}
		data.Response(rw, errRes, http.StatusNotFound)
		return
	}

	res, err := route.client.ValidateEmail(r.Context(), &pb.ValidateEmailParams{
		Token: token,
	})

	if errMsg := common.ConvertRPCError(err); errMsg == common.ErrTokenUsed.Error() {
		errRes := data.FailureResponse{
			Detail: "Token has been used",
			Code:   common.PermissionDenied,
			Status: false,
		}
		data.Response(rw, errRes, http.StatusBadGateway)
		return
	}

	if errMsg := common.ConvertRPCError(err); errMsg == common.ErrExpiredToken.Error() {
		errRes := data.FailureResponse{
			Detail: "Token is expired",
			Code:   common.Unauthorized,
			Status: false,
		}
		data.Response(rw, errRes, http.StatusUnauthorized)
		return
	}

	if isErr := common.ConvertRPCCodeError(err, codes.Internal); isErr {
		errRes := data.FailureResponse{
			Detail: "Internal Server Error",
			Code:   common.InternalError,
			Status: false,
		}
		data.Response(rw, errRes, http.StatusInternalServerError)
		return
	}

	data.Response(rw, res, http.StatusOK)
}
