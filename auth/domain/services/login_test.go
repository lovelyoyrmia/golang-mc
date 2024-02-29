package services

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/Foedie/foedie-server-v2/auth/domain/pb"
	"github.com/Foedie/foedie-server-v2/auth/internal/db"
	"github.com/Foedie/foedie-server-v2/auth/internal/mock"
	"github.com/Foedie/foedie-server-v2/auth/pkg/token"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

type eqLoginUserTxParamsMatcher struct {
	refreshTokenParams db.CreateRefreshTokenParams
	maker              token.Maker
	session            db.Session
}

func (expected eqLoginUserTxParamsMatcher) Matches(x interface{}) bool {
	actualArg, ok := x.(db.CreateRefreshTokenParams)
	if !ok {
		return false
	}

	payload, err := expected.maker.VerifyToken(actualArg.RefreshToken.String)
	if err != nil {
		return false
	}

	return !reflect.DeepEqual(payload.ID.String(), expected.refreshTokenParams.ID)
}

func (e eqLoginUserTxParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and session %v", e.session, e.session)
}

func EqLoginUserTxParams(refreshTokenParams db.CreateRefreshTokenParams, maker token.Maker, session db.Session) gomock.Matcher {
	return eqLoginUserTxParamsMatcher{refreshTokenParams, maker, session}
}

func randomSession(t *testing.T, maker token.Maker, payload *token.Payload, refreshToken string) db.Session {
	user, _ := randomUser(t, true)

	return db.Session{
		ID:  payload.ID.String(),
		Uid: user.Uid,
		RefreshToken: pgtype.Text{
			String: refreshToken,
			Valid:  true,
		},
		ExpiredAt: pgtype.Timestamp{
			Time:  payload.ExpiredAt,
			Valid: true,
		},
		UserAgent: pgtype.Text{
			String: "",
			Valid:  true,
		},
		ClientIp: pgtype.Text{
			String: "",
			Valid:  true,
		},
	}
}

func TestLoginGRPC(t *testing.T) {

	user, password := randomUser(t, true)

	testCases := []struct {
		name          string
		req           *pb.LoginUserParams
		buildStubs    func(store *mock.MockStore, maker token.Maker)
		checkResponse func(t *testing.T, res *pb.LoginUserResponse, err error)
	}{
		{
			name: "OK",
			req: &pb.LoginUserParams{
				Email:    user.Email,
				Password: password,
			},
			buildStubs: func(store *mock.MockStore, maker token.Maker) {
				store.EXPECT().
					GetUserByEmailOrUsername(gomock.Any(), user.Email).
					Times(1).
					Return(user, nil)

				expiredAt := (time.Hour * 24) * 30
				refreshToken, payload, errRefreshToken := maker.GenerateToken(user.Uid, expiredAt)

				require.NoError(t, errRefreshToken)

				refreshTokenParams := db.CreateRefreshTokenParams{
					ID:  payload.ID.String(),
					Uid: user.Uid,
					RefreshToken: pgtype.Text{
						String: refreshToken,
						Valid:  true,
					},
					ExpiredAt: pgtype.Timestamp{
						Time:  payload.ExpiredAt,
						Valid: true,
					},
					UserAgent: pgtype.Text{
						String: "",
						Valid:  true,
					},
					ClientIp: pgtype.Text{
						String: "",
						Valid:  true,
					},
				}

				session := randomSession(t, maker, payload, refreshToken)
				store.EXPECT().
					CreateRefreshToken(gomock.Any(), EqLoginUserTxParams(refreshTokenParams, maker, session)).
					Times(1).
					Return(session, nil)

				store.EXPECT().
					UpdateUserLastLogin(gomock.Any(), user.Uid).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, res *pb.LoginUserResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
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
			tc.buildStubs(store, server.maker)

			res, err := server.LoginUser(context.Background(), tc.req)
			tc.checkResponse(t, res, err)
		})
	}
}
