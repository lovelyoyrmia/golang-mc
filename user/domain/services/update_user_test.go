package services

import (
	"context"
	"testing"

	"github.com/Foedie/foedie-server-v2/user/domain/pb"
	"github.com/Foedie/foedie-server-v2/user/internal/db"
	"github.com/Foedie/foedie-server-v2/user/internal/mock"
	"github.com/Foedie/foedie-server-v2/user/pkg/constants"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestUpdateUserGRPC(t *testing.T) {
	user, _ := randomUser(t, true, true)
	user1, _ := randomUser(t, true, true)

	testCases := []struct {
		name          string
		req           *pb.UserParams
		buildStubs    func(store *mock.MockStore, client *mock.MockAuthServiceClient)
		checkResponse func(t *testing.T, res *pb.SuccessResponse, err error)
	}{
		{
			name: "OK",
			req: &pb.UserParams{
				Uid:         user.Uid,
				Email:       user.Email,
				PhoneNumber: user1.PhoneNumber,
				Username:    user1.Username,
				FirstName:   user1.FirstName,
				LastName:    user1.LastName.String,
			},
			buildStubs: func(store *mock.MockStore, client *mock.MockAuthServiceClient) {

				expectUser := db.GetUserRow{
					Uid:         user.Uid,
					Email:       user.Email,
					Username:    user1.Username,
					FirstName:   user1.FirstName,
					LastName:    user1.LastName,
					PhoneNumber: user1.PhoneNumber,
					IsVerified:  user.IsVerified,
					SecretKey:   user.SecretKey,
					LastLogin:   user.LastLogin,
				}

				store.EXPECT().
					GetUser(gomock.Any(), user.Uid).
					Times(1).
					Return(expectUser, nil)

				store.EXPECT().
					UpdateUserTx(gomock.Any(), expectUser, db.UpdateUserParams{
						Uid:       user.Uid,
						Email:     user.Email,
						Username:  user1.Username,
						FirstName: user1.FirstName,
						LastName: pgtype.Text{
							String: user1.LastName.String,
							Valid:  true,
						},
						PhoneNumber: user1.PhoneNumber,
					}).
					Times(1).
					Return(nil)

			},
			checkResponse: func(t *testing.T, res *pb.SuccessResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
			},
		},
		{
			name: "OK/Email",
			req: &pb.UserParams{
				Uid:         user.Uid,
				Email:       user1.Email,
				PhoneNumber: user.PhoneNumber,
				Username:    user.Username,
				FirstName:   user.FirstName,
				LastName:    user.LastName.String,
			},
			buildStubs: func(store *mock.MockStore, client *mock.MockAuthServiceClient) {

				expectUser := db.GetUserRow{
					Uid:         user.Uid,
					Email:       user1.Email,
					Username:    user.Username,
					FirstName:   user.FirstName,
					LastName:    user.LastName,
					PhoneNumber: user.PhoneNumber,
					IsVerified:  user.IsVerified,
					SecretKey:   user.SecretKey,
					LastLogin:   user.LastLogin,
				}

				store.EXPECT().
					GetUser(gomock.Any(), user.Uid).
					Times(1).
					Return(expectUser, nil)

				store.EXPECT().
					UpdateUserTx(gomock.Any(), expectUser, db.UpdateUserParams{
						Uid:       user.Uid,
						Email:     user1.Email,
						Username:  user.Username,
						FirstName: user.FirstName,
						LastName: pgtype.Text{
							String: user.LastName.String,
							Valid:  true,
						},
						PhoneNumber: user.PhoneNumber,
					}).
					Times(1).
					Return(nil)

				if user.Email != user1.Email {
					client.EXPECT().
						VerifyEmail(user1.Email).
						AnyTimes().
						Return(&pb.VerifyEmailResponse{}, nil)
				}
			},
			checkResponse: func(t *testing.T, res *pb.SuccessResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
			},
		},
		{
			name: "NotFound",
			req: &pb.UserParams{
				Uid:         "",
				Email:       user1.Email,
				PhoneNumber: user.PhoneNumber,
				Username:    user.Username,
				FirstName:   user.FirstName,
				LastName:    user.LastName.String,
			},
			buildStubs: func(store *mock.MockStore, client *mock.MockAuthServiceClient) {

				store.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.GetUserRow{}, constants.ErrRecordNotFound)

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
			req:  &pb.UserParams{},
			buildStubs: func(store *mock.MockStore, client *mock.MockAuthServiceClient) {

				store.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.GetUserRow{}, constants.ErrInternalError)

				store.EXPECT().
					UpdateUserTx(gomock.Any(), gomock.Any(), gomock.Any()).
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
			authClient := mock.NewMockAuthServiceClient(storeCtrl)

			taskCtrl := gomock.NewController(t)
			defer taskCtrl.Finish()

			server := newTestServer(t, store)
			tc.buildStubs(store, authClient)

			res, err := server.UpdateUser(context.Background(), tc.req)
			tc.checkResponse(t, res, err)
		})
	}
}
