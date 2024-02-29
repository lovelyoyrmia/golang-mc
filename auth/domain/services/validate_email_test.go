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

func randomToken(user db.User) string {
	randmString := utils.RandomString(32)
	return fmt.Sprintf("%s__%s", randmString, user.Uid)
}

func randomVerifyEmail(t *testing.T, user db.User, isUsed bool, expiredAt time.Time) db.VerifyEmail {

	token := randomToken(user)

	return db.VerifyEmail{
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

func TestValidateEmailGRPC(t *testing.T) {
	expiredAt := time.Now().Add(time.Minute * 15)

	newExpiredAt := time.Now()

	user, _ := randomUser(t, false)
	newUser, _ := randomUser(t, true)

	token := randomToken(user)

	verifyEmail := randomVerifyEmail(t, user, false, expiredAt)
	newVerifyEmail := randomVerifyEmail(t, newUser, true, newExpiredAt)
	newVerifyEmail1 := randomVerifyEmail(t, user, true, newExpiredAt)

	testCases := []struct {
		name          string
		req           *pb.ValidateEmailParams
		buildStubs    func(store *mock.MockStore, taskUser *mock.MockTaskUser)
		checkResponse func(t *testing.T, res *pb.VerifyEmailResponse, err error)
	}{
		{
			name: "OK",
			req: &pb.ValidateEmailParams{
				Token: verifyEmail.SecretKey,
			},
			buildStubs: func(store *mock.MockStore, taskUser *mock.MockTaskUser) {
				token := verifyEmail.SecretKey
				uid := strings.Split(token, "__")[1]

				store.EXPECT().
					GetVerifyEmail(gomock.Any(), db.GetVerifyEmailParams{
						Uid:       uid,
						SecretKey: token,
					}).
					Times(1).
					Return(verifyEmail, nil)

				store.EXPECT().
					ValidateEmailTx(gomock.Any(), uid, token).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, res *pb.VerifyEmailResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
			},
		},
		{
			name: "TokenUsed",
			req: &pb.ValidateEmailParams{
				Token: newVerifyEmail.SecretKey,
			},
			buildStubs: func(store *mock.MockStore, taskUser *mock.MockTaskUser) {
				token := newVerifyEmail.SecretKey
				uid := strings.Split(token, "__")[1]

				store.EXPECT().
					GetVerifyEmail(gomock.Any(), db.GetVerifyEmailParams{
						Uid:       uid,
						SecretKey: token,
					}).
					Times(1).
					Return(db.VerifyEmail{}, constants.ErrTokenUsed)

				store.EXPECT().
					ValidateEmailTx(gomock.Any(), uid, token).
					AnyTimes().
					Return(constants.ErrInternalError)
			},
			checkResponse: func(t *testing.T, res *pb.VerifyEmailResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Aborted, st.Code())
			},
		},
		{
			name: "NotFound",
			req: &pb.ValidateEmailParams{
				Token: token,
			},
			buildStubs: func(store *mock.MockStore, taskUser *mock.MockTaskUser) {
				uid := strings.Split(token, "__")[1]

				store.EXPECT().
					GetVerifyEmail(gomock.Any(), db.GetVerifyEmailParams{
						Uid:       uid,
						SecretKey: token,
					}).
					Times(1).
					Return(db.VerifyEmail{}, constants.ErrRecordNotFound)

				store.EXPECT().
					ValidateEmailTx(gomock.Any(), uid, token).
					AnyTimes().
					Return(constants.ErrInternalError)
			},
			checkResponse: func(t *testing.T, res *pb.VerifyEmailResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.NotFound, st.Code())
			},
		},
		{
			name: "TokenExpired",
			req: &pb.ValidateEmailParams{
				Token: newVerifyEmail1.SecretKey,
			},
			buildStubs: func(store *mock.MockStore, taskUser *mock.MockTaskUser) {
				token := newVerifyEmail1.SecretKey
				uid := strings.Split(token, "__")[1]

				store.EXPECT().
					GetVerifyEmail(gomock.Any(), db.GetVerifyEmailParams{
						Uid:       uid,
						SecretKey: token,
					}).
					Times(1).
					Return(db.VerifyEmail{}, constants.ErrExpiredToken)

				store.EXPECT().
					ValidateEmailTx(gomock.Any(), uid, token).
					AnyTimes().
					Return(constants.ErrInternalError)
			},
			checkResponse: func(t *testing.T, res *pb.VerifyEmailResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Aborted, st.Code())
			},
		},
		{
			name: "InternalError",
			req: &pb.ValidateEmailParams{
				Token: "test",
			},
			buildStubs: func(store *mock.MockStore, taskUser *mock.MockTaskUser) {
				store.EXPECT().
					ValidateEmailTx(gomock.Any(), gomock.Any(), gomock.Any()).
					AnyTimes().
					Return(constants.ErrInternalError)
			},
			checkResponse: func(t *testing.T, res *pb.VerifyEmailResponse, err error) {
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

			res, err := server.ValidateEmail(context.Background(), tc.req)
			tc.checkResponse(t, res, err)
		})
	}
}
