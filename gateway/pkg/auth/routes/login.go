package routes

import (
	"net/http"

	"github.com/Foedie/foedie-server-v2/gateway/data"
	"github.com/Foedie/foedie-server-v2/gateway/internal/common"
	"github.com/Foedie/foedie-server-v2/gateway/pkg/auth/pb"
	"google.golang.org/grpc/codes"
)

type LoginParams struct {
	Email     string  `json:"email"`
	Password  string  `json:"password"`
	UserAgent *string `json:"-"`
	ClientIP  *string `json:"-"`
}

func (route *AuthRoutes) Login(rw http.ResponseWriter, r *http.Request) {
	loginParams := &LoginParams{}

	err := data.FromJSON(loginParams, r.Body)

	if err != nil {
		errRes := data.FailureResponse{
			Detail: "Unable to decode user params",
			Code:   common.NotValid,
			Status: false,
		}
		data.Response(rw, errRes, http.StatusBadRequest)
		return
	}

	res, err := route.client.LoginUser(r.Context(), &pb.LoginUserParams{
		Email:    loginParams.Email,
		Password: loginParams.Password,
	})

	if err == nil {
		data.Response(rw, res, http.StatusOK)
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

	if errMsg := common.ConvertRPCError(err); errMsg == common.ErrPasswordNotValid.Error() {
		errRes := data.FailureResponse{
			Detail: "Invalid Password",
			Code:   common.Unauthorized,
			Status: false,
		}
		data.Response(rw, errRes, http.StatusUnauthorized)
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

	if errMsg := common.ConvertRPCError(err); errMsg == common.ErrUserInactive.Error() {
		errRes := data.FailureResponse{
			Detail: "User is not active",
			Code:   common.Unauthorized,
			Status: false,
		}
		data.Response(rw, errRes, http.StatusUnauthorized)
		return
	}

	errRes := data.FailureResponse{
		Detail: err.Error(),
		Code:   common.InternalError,
		Status: false,
	}
	data.Response(rw, errRes, http.StatusInternalServerError)
}
