package routes

import (
	"net/http"

	"github.com/Foedie/foedie-server-v2/gateway/data"
	"github.com/Foedie/foedie-server-v2/gateway/internal/common"
	"github.com/Foedie/foedie-server-v2/gateway/internal/models"
	"github.com/Foedie/foedie-server-v2/gateway/pkg/user/pb"
	"google.golang.org/grpc/codes"
)

type UserParams struct {
	Uid         string `json:"uid,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	Email       string `json:"email,omitempty"`
	Username    string `json:"username,omitempty"`
	FirstName   string `json:"first_name,omitempty"`
	LastName    string `json:"last_name,omitempty"`
}

func (route *UserRoutes) UpdateUser(rw http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value(models.KeyUser{}).(string)

	userParams := &UserParams{}

	err := data.FromJSON(userParams, r.Body)

	if err != nil {
		errRes := data.FailureResponse{
			Detail: "Unable to decode user params",
			Code:   common.NotValid,
			Status: false,
		}
		data.Response(rw, errRes, http.StatusBadRequest)
		return
	}

	res, err := route.client.UpdateUser(r.Context(), &pb.UserParams{
		Uid:         uid,
		Email:       userParams.Email,
		PhoneNumber: userParams.PhoneNumber,
		Username:    userParams.Username,
		FirstName:   userParams.FirstName,
		LastName:    userParams.LastName,
	})

	if isErr := common.ConvertRPCCodeError(err, codes.AlreadyExists); isErr {
		errRes := data.FailureResponse{
			Detail: "User already exists",
			Code:   common.InternalError,
			Status: false,
		}
		data.Response(rw, errRes, http.StatusInternalServerError)
		return
	}

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
