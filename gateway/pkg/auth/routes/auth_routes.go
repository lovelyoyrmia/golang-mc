package routes

import (
	"strings"

	"github.com/Foedie/foedie-server-v2/gateway/pkg/auth/pb"
	"github.com/google/uuid"
)

type AuthRoutes struct {
	client pb.AuthServiceClient
}

func NewAuthRoutes(client pb.AuthServiceClient) *AuthRoutes {
	return &AuthRoutes{
		client,
	}
}

func isTokenValid(token string) bool {
	tokenStr := strings.Split(token, "__")

	token = tokenStr[0]
	if len(token) < 32 {
		return false
	}

	uid := uuid.MustParse(tokenStr[1])

	return uid.String() != ""
}
