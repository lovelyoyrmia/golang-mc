package services

import (
	"context"
	"testing"
	"time"

	"github.com/Foedie/foedie-server-v2/auth/domain/pb"
	"github.com/Foedie/foedie-server-v2/auth/internal/db"
	"github.com/Foedie/foedie-server-v2/auth/internal/mock"
	"github.com/Foedie/foedie-server-v2/auth/pkg/constants"
	"github.com/Foedie/foedie-server-v2/auth/pkg/token"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func randomJwtToken(maker token.Maker, user db.User) (string, *token.Payload) {
	token, payload, err := maker.GenerateToken(user.Uid, time.Second*2)
	if err != nil {
		return "", nil
	}
	return token, payload
}

func randomPayload(uid string) *token.Payload {
	return &token.Payload{
		ID:        uuid.New(),
		UID:       uid,
		ExpiredAt: time.Now().Add(time.Second * 2),
	}
}

func TestValidateTokenGRPC(t *testing.T) {
	user, _ := randomUser(t, true)
	payload := randomPayload(user.Uid)

	testCases := []struct {
		name          string
		req           func(string) *pb.ValidateTokenParams
		buildAuth     func(maker token.Maker, user db.User) (string, *token.Payload)
		buildStubs    func(store *mock.MockStore, maker *mock.MockMaker, token string)
		checkResponse func(t *testing.T, res *pb.ValidateTokenResponse, err error)
	}{
		{
			name: "OK",
			req: func(s string) *pb.ValidateTokenParams {
				return &pb.ValidateTokenParams{
					Token: s,
				}
			},
			buildAuth: func(maker token.Maker, user db.User) (string, *token.Payload) {
				return randomJwtToken(maker, user)
			},
			buildStubs: func(store *mock.MockStore, maker *mock.MockMaker, token string) {

				maker.EXPECT().
					VerifyToken(token).
					AnyTimes().
					Return(payload, nil)

				store.EXPECT().
					GetUserByUid(gomock.Any(), user.Uid).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, res *pb.ValidateTokenResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
			},
		},
		{
			name: "InternalError",
			req: func(s string) *pb.ValidateTokenParams {
				return &pb.ValidateTokenParams{
					Token: s,
				}
			},
			buildAuth: func(maker token.Maker, user db.User) (string, *token.Payload) {
				return randomJwtToken(maker, user)
			},
			buildStubs: func(store *mock.MockStore, maker *mock.MockMaker, jwt string) {

				maker.EXPECT().
					VerifyToken("").
					AnyTimes().
					Return(nil, constants.ErrInternalError)

				store.EXPECT().
					GetUserByUid(gomock.Any(), gomock.Any()).
					AnyTimes().
					Return(user, constants.ErrInternalError)
			},
			checkResponse: func(t *testing.T, res *pb.ValidateTokenResponse, err error) {
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
			maker := mock.NewMockMaker(taskCtrl)

			server := newTestServer(t, store, taskUser)
			jwt, _ := tc.buildAuth(server.maker, user)

			req := tc.req(jwt)
			tc.buildStubs(store, maker, jwt)

			res, err := server.ValidateToken(context.Background(), req)
			tc.checkResponse(t, res, err)
		})
	}
}
