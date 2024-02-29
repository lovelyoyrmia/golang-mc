package routes

import (
	"github.com/Foedie/foedie-server-v2/gateway/pkg/user/pb"
)

type UserRoutes struct {
	client pb.UserServiceClient
}

func NewUserRoutes(client pb.UserServiceClient) *UserRoutes {
	return &UserRoutes{
		client,
	}
}
