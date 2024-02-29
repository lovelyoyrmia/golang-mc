package routes

import (
	"net/http"

	"github.com/Foedie/foedie-server-v2/gateway/data"
	"github.com/Foedie/foedie-server-v2/gateway/internal/common"
	"github.com/Foedie/foedie-server-v2/gateway/pkg/auth/pb"
	"google.golang.org/grpc/codes"
)

type VerifyEmailParams struct {
	Email string `json:"email"`
}

func (route *AuthRoutes) VerifyEmail(rw http.ResponseWriter, r *http.Request) {
	verifyEmailParams := &VerifyEmailParams{}

	err := data.FromJSON(verifyEmailParams, r.Body)

	if err != nil {
		errRes := data.FailureResponse{
			Detail: "Unable to decode user params",
			Code:   common.NotValid,
			Status: false,
		}
		data.Response(rw, errRes, http.StatusBadRequest)
		return
	}

	res, err := route.client.VerifyEmail(r.Context(), &pb.VerifyEmailParams{
		Email: verifyEmailParams.Email,
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

	if isErr := common.ConvertRPCCodeError(err, codes.Aborted); isErr {
		errRes := data.FailureResponse{
			Detail: "User has been verified",
			Code:   common.PermissionDenied,
			Status: false,
		}
		data.Response(rw, errRes, http.StatusBadGateway)
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
