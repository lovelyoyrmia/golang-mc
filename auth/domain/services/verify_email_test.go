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

func TestVerifyEmailGRPC(t *testing.T) {
	user, _ := randomUser(t, false)
	newUser, _ := randomUser(t, true)

	testCases := []struct {
		name          string
		req           string
		buildStubs    func(store *mock.MockStore, taskUser *mock.MockTaskUser)
		checkResponse func(t *testing.T, res *pb.VerifyEmailResponse, err error)
	}{
		{
			name: "OK",
			req:  user.Email,
			buildStubs: func(store *mock.MockStore, taskUser *mock.MockTaskUser) {
				store.EXPECT().
					GetUserByEmailOrUsername(gomock.Any(), user.Email).
					Times(1).
					Return(user, nil)
				taskPayload := &worker.PayloadSendVerifyEmail{
					Email: user.Email,
				}
				taskUser.EXPECT().
					UserTaskSendVerificationEmail(gomock.Any(), taskPayload, gomock.Any()).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, res *pb.VerifyEmailResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
			},
		},
		{
			name: "InternalError",
			req:  user.Email,
			buildStubs: func(store *mock.MockStore, taskUser *mock.MockTaskUser) {
				store.EXPECT().GetUserByEmailOrUsername(gomock.Any(), gomock.Any()).
					Times(1).Return(db.User{}, constants.ErrInternalError)

				taskUser.EXPECT().
					UserTaskSendVerificationEmail(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(constants.ErrInternalError).AnyTimes()
			},
			checkResponse: func(t *testing.T, res *pb.VerifyEmailResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Internal, st.Code())
			},
		},
		{
			name: "NotFound",
			req:  "test@gmail.com",
			buildStubs: func(store *mock.MockStore, taskUser *mock.MockTaskUser) {
				store.EXPECT().GetUserByEmailOrUsername(gomock.Any(), "test@gmail.com").
					Times(1).Return(db.User{}, constants.ErrRecordNotFound)

				taskUser.EXPECT().
					UserTaskSendVerificationEmail(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(constants.ErrInternalError).AnyTimes()
			},
			checkResponse: func(t *testing.T, res *pb.VerifyEmailResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.NotFound, st.Code())
			},
		},
		{
			name: "UserVerified",
			req:  newUser.Email,
			buildStubs: func(store *mock.MockStore, taskUser *mock.MockTaskUser) {
				store.EXPECT().GetUserByEmailOrUsername(gomock.Any(), newUser.Email).
					Times(1).Return(db.User{}, constants.ErrUserVerified)

				taskUser.EXPECT().
					UserTaskSendVerificationEmail(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(constants.ErrInternalError).AnyTimes()
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

			res, err := server.VerifyEmail(context.Background(), &pb.VerifyEmailParams{
				Email: tc.req,
			})
			tc.checkResponse(t, res, err)
		})
	}
}
