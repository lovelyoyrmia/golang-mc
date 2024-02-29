package routes

import (
	"net/http"

	"github.com/Foedie/foedie-server-v2/gateway/data"
	"github.com/Foedie/foedie-server-v2/gateway/internal/common"
	"github.com/Foedie/foedie-server-v2/gateway/internal/models"
	"github.com/Foedie/foedie-server-v2/gateway/pkg/user/pb"
	"google.golang.org/grpc/codes"
)

func (route *UserRoutes) DeleteUser(rw http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value(models.KeyUser{}).(string)

	if uid == "" {
		errRes := data.FailureResponse{
			Detail: "Token is required",
			Code:   common.Unauthorized,
			Status: false,
		}
		data.Response(rw, errRes, http.StatusUnauthorized)
		return
	}

	res, err := route.client.DeleteUser(r.Context(), &pb.UserUidParams{
		Uid: uid,
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

	if err != nil {
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
