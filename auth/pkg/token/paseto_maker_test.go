package token

import (
	"testing"
	"time"

	"github.com/Foedie/foedie-server-v2/auth/pkg/config"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {
	conf := config.Config{
		SecretKey: "imf580tLVqWu9RLqCvFNCuHyt21zasca",
	}
	maker, err := NewPasetoMaker(conf)
	require.NoError(t, err)

	uid := uuid.New().String()
	duration := time.Minute

	token, _, err := maker.GenerateToken(uid, duration)
	require.NoError(t, err)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)

	require.Equal(t, payload.UID, uid)
}
