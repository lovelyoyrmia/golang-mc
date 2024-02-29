package services

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/Foedie/foedie-server-v2/auth/domain/pb"
	"github.com/Foedie/foedie-server-v2/auth/internal/db"
	"github.com/Foedie/foedie-server-v2/auth/internal/mock"
	"github.com/Foedie/foedie-server-v2/auth/pkg/constants"
	"github.com/Foedie/foedie-server-v2/auth/pkg/utils"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func randomRecoverAccount(t *testing.T, user db.User, isUsed bool, expiredAt time.Time) db.RecoverAccount {

	randmString := utils.RandomString(32)
	token := fmt.Sprintf("%s__%s", randmString, user.Uid)

	return db.RecoverAccount{
		ID:        uuid.NewString(),
		Uid:       user.Uid,
		Email:     user.Email,
		SecretKey: token,
		IsUsed:    isUsed,
		ExpiredAt: pgtype.Timestamp{
			Time:  expiredAt,
			Valid: true,
		},
	}
}

func TestValidateRecoverAccountGRPC(t *testing.T) {
	expiredAt := time.Now().Add(time.Minute * 15)

	newExpiredAt := time.Now()

	user, _ := randomUser(t, false)
	newUser, _ := randomUser(t, true)

	token := randomToken(user)

	recoverAccount := randomRecoverAccount(t, user, false, expiredAt)
	newRecoverAccount := randomRecoverAccount(t, newUser, true, newExpiredAt)
	newRecoverAccount1 := randomRecoverAccount(t, user, true, newExpiredAt)

	testCases := []struct {
		name          string
		req           *pb.ValidateRecoverAccountParams
		buildStubs    func(store *mock.MockStore, taskUser *mock.MockTaskUser)
		checkResponse func(t *testing.T, res *pb.RecoverAccountResponse, err error)
	}{
		{
			name: "OK",
			req: &pb.ValidateRecoverAccountParams{
				Token: recoverAccount.SecretKey,
			},
			buildStubs: func(store *mock.MockStore, taskUser *mock.MockTaskUser) {
				newToken := recoverAccount.SecretKey
				uid := strings.Split(newToken, "__")[1]

				store.EXPECT().
					GetRecoverAccount(gomock.Any(), db.GetRecoverAccountParams{
						Uid:       uid,
						SecretKey: newToken,
					}).
					Times(1).
					Return(recoverAccount, nil)

				store.EXPECT().
					ValidateRecoverAccountTx(gomock.Any(), uid, newToken).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, res *pb.RecoverAccountResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
			},
		},
		{
			name: "TokenUsed",
			req: &pb.ValidateRecoverAccountParams{
				Token: newRecoverAccount.SecretKey,
			},
			buildStubs: func(store *mock.MockStore, taskUser *mock.MockTaskUser) {
				token := newRecoverAccount.SecretKey
				uid := strings.Split(token, "__")[1]

				store.EXPECT().
					GetRecoverAccount(gomock.Any(), db.GetRecoverAccountParams{
						Uid:       uid,
						SecretKey: token,
					}).
					Times(1).
					Return(db.RecoverAccount{}, constants.ErrTokenUsed)

				store.EXPECT().
					ValidateRecoverAccountTx(gomock.Any(), uid, token).
					AnyTimes().
					Return(constants.ErrInternalError)
			},
			checkResponse: func(t *testing.T, res *pb.RecoverAccountResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Aborted, st.Code())
			},
		},
		{
			name: "TokenExpired",
			req: &pb.ValidateRecoverAccountParams{
				Token: newRecoverAccount1.SecretKey,
			},
			buildStubs: func(store *mock.MockStore, taskUser *mock.MockTaskUser) {
				token := newRecoverAccount1.SecretKey
				uid := strings.Split(token, "__")[1]

				store.EXPECT().
					GetRecoverAccount(gomock.Any(), db.GetRecoverAccountParams{
						Uid:       uid,
						SecretKey: token,
					}).
					Times(1).
					Return(db.RecoverAccount{}, constants.ErrExpiredToken)

				store.EXPECT().
					ValidateRecoverAccountTx(gomock.Any(), uid, token).
					AnyTimes().
					Return(constants.ErrInternalError)
			},
			checkResponse: func(t *testing.T, res *pb.RecoverAccountResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Aborted, st.Code())
			},
		},
		{
			name: "NotFound",
			req: &pb.ValidateRecoverAccountParams{
				Token: token,
			},
			buildStubs: func(store *mock.MockStore, taskUser *mock.MockTaskUser) {
				uid := strings.Split(token, "__")[1]

				store.EXPECT().
					GetRecoverAccount(gomock.Any(), db.GetRecoverAccountParams{
						Uid:       uid,
						SecretKey: token,
					}).
					Times(1).
					Return(db.RecoverAccount{}, constants.ErrRecordNotFound)

				store.EXPECT().
					ValidateRecoverAccountTx(gomock.Any(), uid, token).
					AnyTimes().
					Return(constants.ErrInternalError)
			},
			checkResponse: func(t *testing.T, res *pb.RecoverAccountResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.NotFound, st.Code())
			},
		},
		{
			name: "InternalError",
			req: &pb.ValidateRecoverAccountParams{
				Token: "test",
			},
			buildStubs: func(store *mock.MockStore, taskUser *mock.MockTaskUser) {
				store.EXPECT().
					ValidateRecoverAccountTx(gomock.Any(), gomock.Any(), gomock.Any()).
					AnyTimes().
					Return(constants.ErrInternalError)
			},
			checkResponse: func(t *testing.T, res *pb.RecoverAccountResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Internal, st.Code())
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			storeCtrl := gomock.NewController(t)
			defer storeCtrl.Finish()
			store := mock.NewMockStore(storeCtrl)

			taskCtrl := gomock.NewController(t)
			defer taskCtrl.Finish()

			taskUser := mock.NewMockTaskUser(taskCtrl)

			server := newTestServer(t, store, taskUser)
			tc.buildStubs(store, taskUser)

			res, err := server.ValidateRecoverAccount(context.Background(), tc.req)
			tc.checkResponse(t, res, err)
		})
	}
}
