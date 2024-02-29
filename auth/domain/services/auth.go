package services

import (
	"context"
	"strings"

	"github.com/Foedie/foedie-server-v2/auth/domain/pb"
	"github.com/Foedie/foedie-server-v2/auth/internal/db"
	"github.com/Foedie/foedie-server-v2/auth/internal/worker"
	"github.com/Foedie/foedie-server-v2/auth/pkg/constants"
	"github.com/Foedie/foedie-server-v2/auth/pkg/token"
	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	grpcUserAgentHeader string = "grpcgateway-user-agent"
	userAgentHeader     string = "user-agent"
	xForwardedForHeader string = "x-forwarded-for"
)

type Metadata struct {
	UserAgent string
	ClientIP  string
}

type Server struct {
	pb.UnimplementedAuthServiceServer
	store    db.Store
	maker    token.Maker
	taskUser worker.TaskUser
}

func NewServer(store db.Store, maker token.Maker, taskUser worker.TaskUser) *Server {
	return &Server{
		store:    store,
		maker:    maker,
		taskUser: taskUser,
	}
}

func extractMetadata(ctx context.Context) *Metadata {
	m := &Metadata{}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if userAgents := md.Get(grpcUserAgentHeader); len(userAgents) > 0 {
			m.UserAgent = userAgents[0]
		}
		if grpcUserAgents := md.Get(userAgentHeader); len(grpcUserAgents) > 0 {
			m.UserAgent = grpcUserAgents[0]
		}
		if clientIP := md.Get(xForwardedForHeader); len(clientIP) > 0 {
			m.ClientIP = clientIP[0]
		}
	}

	if p, ok := peer.FromContext(ctx); ok {
		m.ClientIP = p.Addr.String()
	}

	return m
}

func tokenValid(token string) (string, error) {
	tokenStr := strings.Split(token, "__")

	token = tokenStr[0]
	if len(token) < 32 {
		return "", constants.ErrInternalError
	}

	uid, err := uuid.Parse(tokenStr[1])

	if err != nil {
		return "", err
	}

	return uid.String(), nil
}
