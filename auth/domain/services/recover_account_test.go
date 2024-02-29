package services

import (
	"context"
	"testing"

	"github.com/Foedie/foedie-server-v2/auth/domain/pb"
	"github.com/Foedie/foedie-server-v2/auth/internal/db"
	"github.com/Foedie/foedie-server-v2/auth/internal/mock"
	"github.com/Foedie/foedie-server-v2/auth/internal/worker"
	"github.com/Foedie/foedie-server-v2/auth/pkg/constants"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestRecoverAccountGRPC(t *testing.T) {
	user, _ := randomUser(t, true, false)
	newUserActive, _ := randomUser(t, true, true)
	newUserVerified, _ := randomUser(t, false, true)

	testCases := []struct {
		name          string
		req           *pb.RecoverAccountParams
		buildStubs    func(store *mock.MockStore, taskUser *mock.MockTaskUser)
		checkResponse func(t *testing.T, res *pb.RecoverAccountResponse, err error)
	}{
		{
			name: "OK",
			req: &pb.RecoverAccountParams{
				Email:    user.Email,
				Username: user.Username,
			},
			buildStubs: func(store *mock.MockStore, taskUser *mock.MockTaskUser) {
				store.EXPECT().
					GetUserByEmailAndUsername(gomock.Any(), db.GetUserByEmailAndUsernameParams{
						Username: user.Username,
						Email:    user.Email,
					}).
					Times(1).
					Return(user, nil)
				taskPayload := &worker.PayloadSendRecoverAccount{
					Email:    user.Email,
					Username: user.Username,
				}
				taskUser.EXPECT().
					UserTaskSendRecoverAccount(gomock.Any(), taskPayload, gomock.Any()).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, res *pb.RecoverAccountResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
			},
		},
		{
			name: "NotFound",
			req: &pb.RecoverAccountParams{
				Email:    user.Email,
				Username: user.Username,
			},
			buildStubs: func(store *mock.MockStore, taskUser *mock.MockTaskUser) {
				store.EXPECT().
					GetUserByEmailAndUsername(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, constants.ErrRecordNotFound)

				taskUser.EXPECT().
					UserTaskSendRecoverAccount(gomock.Any(), gomock.Any(), gomock.Any()).
					AnyTimes().Return(constants.ErrInternalError)
			},
			checkResponse: func(t *testing.T, res *pb.RecoverAccountResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.NotFound, st.Code())
			},
		},
		{
			name: "UserActive",
			req: &pb.RecoverAccountParams{
				Email:    newUserActive.Email,
				Username: newUserActive.Username,
			},
			buildStubs: func(store *mock.MockStore, taskUser *mock.MockTaskUser) {
				store.EXPECT().
					GetUserByEmailAndUsername(gomock.Any(), db.GetUserByEmailAndUsernameParams{
						Email:    newUserActive.Email,
						Username: newUserActive.Username,
					}).
					Times(1).
					Return(db.User{}, constants.ErrUserActive)

				taskUser.EXPECT().
					UserTaskSendRecoverAccount(gomock.Any(), gomock.Any(), gomock.Any()).
					AnyTimes().Return(constants.ErrInternalError)
			},
			checkResponse: func(t *testing.T, res *pb.RecoverAccountResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Aborted, st.Code())
			},
		},
		{
			name: "NotVerified",
			req: &pb.RecoverAccountParams{
				Email:    newUserVerified.Email,
				Username: newUserVerified.Username,
			},
			buildStubs: func(store *mock.MockStore, taskUser *mock.MockTaskUser) {
				store.EXPECT().
					GetUserByEmailAndUsername(gomock.Any(), db.GetUserByEmailAndUsernameParams{
						Email:    newUserVerified.Email,
						Username: newUserVerified.Username,
					}).
					Times(1).
					Return(db.User{}, constants.ErrUserNotVerified)

				taskUser.EXPECT().
					UserTaskSendRecoverAccount(gomock.Any(), gomock.Any(), gomock.Any()).
					AnyTimes().Return(constants.ErrInternalError)
			},
			checkResponse: func(t *testing.T, res *pb.RecoverAccountResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Aborted, st.Code())
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

			res, err := server.RecoverAccount(context.Background(), tc.req)
			tc.checkResponse(t, res, err)
		})
	}
}
