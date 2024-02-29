package services

import (
	"github.com/Foedie/foedie-server-v2/email/domain/pb"
	"github.com/Foedie/foedie-server-v2/email/pkg/config"
)

type Server struct {
	pb.UnimplementedEmailServiceServer
	config config.Config
}

func NewServer(config config.Config) *Server {
	return &Server{
		config: config,
	}
}
