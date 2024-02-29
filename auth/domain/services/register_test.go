package services

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/Foedie/foedie-server-v2/auth/domain/pb"
	"github.com/Foedie/foedie-server-v2/auth/internal/db"
	"github.com/Foedie/foedie-server-v2/auth/internal/mock"
	"github.com/Foedie/foedie-server-v2/auth/internal/worker"
	"github.com/Foedie/foedie-server-v2/auth/pkg/constants"
	"github.com/Foedie/foedie-server-v2/auth/pkg/utils"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type eqCreateUserTxParamsMatcher struct {
	arg      db.UserTxParams
	password string
	user     db.User
}

func (expected eqCreateUserTxParamsMatcher) Matches(x interface{}) bool {
	actualArg, ok := x.(db.UserTxParams)
	if !ok {
		return false
	}

	err := utils.ComparePassword(actualArg.UserParams.Password, expected.password)
	if err != nil {
		return false
	}

	decrypt, err := utils.Decrypt(expected.user.SecretKey)
	if err != nil {
		return false
	}

	if !reflect.DeepEqual(decrypt, actualArg.UserParams.Uid) {
		return false
	}

	expected.arg.UserParams.Password = actualArg.UserParams.Password
	if !reflect.DeepEqual(expected.arg.UserParams, actualArg.UserParams) {
		return false
	}

	err = actualArg.AfterCreate(expected.user)
	return err == nil
}

func (e eqCreateUserTxParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqCreateUserTxParams(arg db.UserTxParams, password string, user db.User) gomock.Matcher {
	return eqCreateUserTxParamsMatcher{arg, password, user}
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

func TestCreateUserGRPC(t *testing.T) {
	user, password := randomUser(t, false)

	testCases := []struct {
		name          string
		req           *pb.CreateUserParams
		buildStubs    func(store *mock.MockStore, taskUser *mock.MockTaskUser)
		checkResponse func(t *testing.T, res *pb.CreateUserResponse, err error)
	}{
		{
			name: "OK",
			req: &pb.CreateUserParams{
				Uid:             user.Uid,
				Username:        user.Username,
				Password:        password,
				PhoneNumber:     user.PhoneNumber,
				FirstName:       user.FirstName,
				LastName:        user.LastName.String,
				Email:           user.Email,
				ConfirmPassword: password,
			},
			buildStubs: func(store *mock.MockStore, taskUser *mock.MockTaskUser) {
				arg := db.UserTxParams{
					UserParams: db.CreateUserParams{
						Uid:         user.Uid,
						Email:       user.Email,
						PhoneNumber: user.PhoneNumber,
						Username:    user.Username,
						FirstName:   user.FirstName,
						LastName:    user.LastName,
						SecretKey:   user.SecretKey,
						IsActive:    user.IsActive,
						IsVerified:  user.IsVerified,
					},
				}

				store.EXPECT().
					CreateUserTx(gomock.Any(), EqCreateUserTxParams(arg, password, user)).
					Times(1).
					Return(nil)

				taskPayload := &worker.PayloadSendVerifyEmail{
					Email: user.Email,
				}
				taskUser.EXPECT().
					UserTaskSendVerificationEmail(gomock.Any(), taskPayload, gomock.Any()).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, res *pb.CreateUserResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
			},
		},
		{
			name: "InternalError",
			req: &pb.CreateUserParams{
				Uid:             user.Uid,
				Username:        user.Username,
				Password:        password,
				PhoneNumber:     user.PhoneNumber,
				FirstName:       user.FirstName,
				LastName:        user.LastName.String,
				Email:           user.Email,
				ConfirmPassword: password,
			},
			buildStubs: func(store *mock.MockStore, taskUser *mock.MockTaskUser) {
				store.EXPECT().CreateUserTx(gomock.Any(), gomock.Any()).
					Times(1).Return(constants.ErrInternalError)

				taskUser.EXPECT().
					UserTaskSendVerificationEmail(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *pb.CreateUserResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Internal, st.Code())
			},
		},
		{
			name: "AlreadyExists",
			req: &pb.CreateUserParams{
				Uid:             user.Uid,
				Username:        user.Username,
				Password:        password,
				PhoneNumber:     user.PhoneNumber,
				FirstName:       user.FirstName,
				LastName:        user.LastName.String,
				Email:           user.Email,
				ConfirmPassword: password,
			},
			buildStubs: func(store *mock.MockStore, taskDistributor *mock.MockTaskUser) {
				store.EXPECT().
					CreateUserTx(gomock.Any(), gomock.Any()).
					Times(1).
					Return(constants.ErrUniqueViolation)

				taskDistributor.EXPECT().
					UserTaskSendVerificationEmail(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *pb.CreateUserResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.AlreadyExists, st.Code())
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

			tc.buildStubs(store, taskUser)
			server := newTestServer(t, store, taskUser)

			res, err := server.CreateUser(context.Background(), tc.req)
			tc.checkResponse(t, res, err)
		})
	}
}
