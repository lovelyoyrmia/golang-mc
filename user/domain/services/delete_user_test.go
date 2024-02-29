package services

import (
	"context"
	"testing"

	"github.com/Foedie/foedie-server-v2/user/domain/pb"
	"github.com/Foedie/foedie-server-v2/user/internal/mock"
	"github.com/Foedie/foedie-server-v2/user/pkg/constants"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestDeleteUserGRPC(t *testing.T) {
	user, _ := randomUser(t, true, true)

	testCases := []struct {
		name          string
		req           *pb.UserUidParams
		buildStubs    func(store *mock.MockStore)
		checkResponse func(t *testing.T, res *pb.SuccessResponse, err error)
	}{
		{
			name: "OK",
			req: &pb.UserUidParams{
				Uid: user.Uid,
			},
			buildStubs: func(store *mock.MockStore) {

				store.EXPECT().
					DeleteUser(gomock.Any(), user.Uid).
					Times(1).
					Return(nil)

			},
			checkResponse: func(t *testing.T, res *pb.SuccessResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
			},
		},
		{
			name: "NotFound",
			req: &pb.UserUidParams{
				Uid: "",
			},
			buildStubs: func(store *mock.MockStore) {

				store.EXPECT().
					DeleteUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(constants.ErrRecordNotFound)

			},
			checkResponse: func(t *testing.T, res *pb.SuccessResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.NotFound, st.Code())
			},
		},
		{
			name: "InternalError",
			req:  &pb.UserUidParams{},
			buildStubs: func(store *mock.MockStore) {

				store.EXPECT().
					DeleteUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(constants.ErrInternalError)

			},
			checkResponse: func(t *testing.T, res *pb.SuccessResponse, err error) {
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

			server := newTestServer(t, store)
			tc.buildStubs(store)

			res, err := server.DeleteUser(context.Background(), tc.req)
			tc.checkResponse(t, res, err)
		})
	}
}
