package services

import (
	"testing"

	"github.com/Foedie/foedie-server-v2/auth/internal/db"
	"github.com/Foedie/foedie-server-v2/auth/internal/worker"
	"github.com/Foedie/foedie-server-v2/auth/pkg/config"
	"github.com/Foedie/foedie-server-v2/auth/pkg/token"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store db.Store, taskUser worker.TaskUser) *Server {
	config := config.Config{
		SecretKey: "imf580tLVqWu9RLqCvFNCuHQN33YhztZ",
	}

	tokenMaker, err := token.NewPasetoMaker(config)
	require.NoError(t, err)

	server := NewServer(store, tokenMaker, taskUser)

	return server
}
