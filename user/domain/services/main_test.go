package services

import (
	"testing"

	"github.com/Foedie/foedie-server-v2/auth/pkg/utils"
	"github.com/Foedie/foedie-server-v2/user/domain/clients"
	"github.com/Foedie/foedie-server-v2/user/internal/db"
	"github.com/Foedie/foedie-server-v2/user/pkg/config"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	config := config.Config{
		SecretKey:  "imf580tLVqWu9RLqCvFNCuHQN33YhztZ",
		AuthSvcUrl: "localhost:50051",
	}

	auth := clients.InitAuthServiceClient(config.AuthSvcUrl)

	server := NewServer(store, auth)

	return server
}

func randomUser(t *testing.T, isVerified bool, IsActive ...bool) (user db.User, password string) {
	password = utils.RandomString(6)
	hashedPassword, err := utils.HashPassword(password)
	require.NoError(t, err)

	uid := uuid.NewString()
	secretKey, err := utils.Encrypt(uid)

	require.NoError(t, err)

	var isActive bool
	if len(IsActive) == 0 {
		isActive = true
	} else {
		isActive = IsActive[0]
	}

	user = db.User{
		Username:  utils.RandomString(15),
		Uid:       uid,
		Email:     utils.GenerateRandomEmail(15),
		FirstName: utils.RandomString(15),
		LastName: pgtype.Text{
			String: utils.RandomString(15),
			Valid:  true,
		},
		Password:    hashedPassword,
		PhoneNumber: utils.RandomString(10),
		SecretKey:   secretKey,
		IsActive: pgtype.Bool{
			Bool:  isActive,
			Valid: true,
		},
		IsVerified: pgtype.Bool{
			Bool:  isVerified,
			Valid: true,
		},
	}
	return
}
