package routes

import (
	"net/http"

	"github.com/Foedie/foedie-server-v2/gateway/data"
	"github.com/Foedie/foedie-server-v2/gateway/internal/common"
	"github.com/Foedie/foedie-server-v2/gateway/pkg/auth/pb"
	"google.golang.org/grpc/codes"
)

type RecoverAccountParams struct {
	Email    string `json:"email"`
	Username string `json:"username"`
}

func (route *AuthRoutes) RecoverAccount(rw http.ResponseWriter, r *http.Request) {
	recAcc := &RecoverAccountParams{}

	err := data.FromJSON(recAcc, r.Body)

	if err != nil {
		errRes := data.FailureResponse{
			Detail: "Unable to decode user params",
			Code:   common.NotValid,
			Status: false,
		}
		data.Response(rw, errRes, http.StatusBadRequest)
		return
	}

	res, err := route.client.RecoverAccount(r.Context(), &pb.RecoverAccountParams{
		Email:    recAcc.Email,
		Username: recAcc.Username,
	})

	if isErr := common.ConvertRPCCodeError(err, codes.NotFound); isErr {
		errRes := data.FailureResponse{
			Detail: "User not found",
			Code:   common.NotFound,
			Status: false,
		}
		data.Response(rw, errRes, http.StatusNotFound)
		return
	}

	if errMsg := common.ConvertRPCError(err); errMsg == common.ErrUserActive.Error() {
		errRes := data.FailureResponse{
			Detail: "User is active",
			Code:   common.PermissionDenied,
			Status: false,
		}
		data.Response(rw, errRes, http.StatusBadGateway)
		return
	}

	if errMsg := common.ConvertRPCError(err); errMsg == common.ErrUserNotVerified.Error() {
		errRes := data.FailureResponse{
			Detail: "User is not verified",
			Code:   common.Unauthorized,
			Status: false,
		}
		data.Response(rw, errRes, http.StatusUnauthorized)
		return
	}

	if isErr := common.ConvertRPCCodeError(err, codes.Internal); isErr {
		errRes := data.FailureResponse{
			Detail: err.Error(),
			Code:   common.InternalError,
			Status: false,
		}
		data.Response(rw, errRes, http.StatusInternalServerError)
		return
	}

	data.Response(rw, res, http.StatusOK)
}
