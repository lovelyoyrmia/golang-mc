package services

import (
	"github.com/Foedie/foedie-server-v2/user/domain/clients"
	"github.com/Foedie/foedie-server-v2/user/domain/pb"
	"github.com/Foedie/foedie-server-v2/user/internal/db"
)

type Server struct {
	pb.UnimplementedUserServiceServer
	store   db.Store
	authSvc clients.AuthServiceClient
}

func NewServer(store db.Store, authSvc clients.AuthServiceClient) *Server {
	return &Server{
		store:   store,
		authSvc: authSvc,
	}
}
