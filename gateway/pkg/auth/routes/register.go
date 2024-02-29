package routes

import (
	"net/http"

	"github.com/Foedie/foedie-server-v2/gateway/data"
	"github.com/Foedie/foedie-server-v2/gateway/internal/common"
	"github.com/Foedie/foedie-server-v2/gateway/pkg/auth/pb"
	"google.golang.org/grpc/codes"
)

type UserParams struct {
	ID              string `json:"-"`
	Username        string `json:"username"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Email           string `json:"email"`
	PhoneNumber     string `json:"phone_number"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

func (route *AuthRoutes) Register(rw http.ResponseWriter, r *http.Request) {
	uParams := &UserParams{}

	err := data.FromJSON(uParams, r.Body)

	if err != nil {
		errRes := data.FailureResponse{
			Detail: "Unable to decode user params",
			Code:   common.NotValid,
			Status: false,
		}
		data.Response(rw, errRes, http.StatusBadRequest)
		return
	}

	// Returns an error if password and confirm password does not match
	if uParams.Password != uParams.ConfirmPassword {
		errRes := data.FailureResponse{
			Detail: "Password did not match",
			Code:   common.NotValid,
			Status: false,
		}
		data.Response(rw, errRes, http.StatusBadRequest)
		return
	}

	res, err := route.client.CreateUser(r.Context(), &pb.CreateUserParams{
		PhoneNumber:     uParams.PhoneNumber,
		Username:        uParams.Username,
		Email:           uParams.Email,
		FirstName:       uParams.FirstName,
		LastName:        uParams.LastName,
		Password:        uParams.Password,
		ConfirmPassword: uParams.ConfirmPassword,
	})

	if err == nil {
		data.Response(rw, res, http.StatusOK)
		return
	}

	if isErr := common.ConvertRPCCodeError(err, codes.AlreadyExists); isErr {
		errRes := data.FailureResponse{
			Detail: "User already exists",
			Code:   common.InternalError,
			Status: false,
		}
		data.Response(rw, errRes, http.StatusInternalServerError)
		return
	}

	errRes := data.FailureResponse{
		Detail: err.Error(),
		Code:   common.InternalError,
		Status: false,
	}
	data.Response(rw, errRes, http.StatusInternalServerError)
}
