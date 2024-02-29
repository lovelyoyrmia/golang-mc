package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

func (store *SQLStore) CreateUserTx(ctx context.Context, req UserTxParams) (err error) {

	err = store.ExecTx(ctx, func(q *Queries) error {

		user, errRes := q.CreateUser(ctx, req.UserParams)
		if errRes != nil {
			return errRes
		}

		err := q.CreateUserOTP(ctx, CreateUserOTPParams{
			Uid: user.Uid,
		})

		if err != nil {
			return err
		}

		return req.AfterCreate(user)
	})
	return
}

func (store *SQLStore) ValidateEmailTx(ctx context.Context, uid string, token string) (err error) {
	err = store.ExecTx(ctx, func(q *Queries) error {
		if err := q.UpdateVerifyEmail(ctx, UpdateVerifyEmailParams{
			Uid:       uid,
			SecretKey: token,
		}); err != nil {
			return err
		}
		return q.UpdateUserVerified(ctx, UpdateUserVerifiedParams{
			Uid: uid,
			IsVerified: pgtype.Bool{
				Bool:  true,
				Valid: true,
			},
		})
	})
	return
}

func (store *SQLStore) ValidateRecoverAccountTx(ctx context.Context, uid string, token string) (err error) {
	err = store.ExecTx(ctx, func(q *Queries) error {
		err = q.UpdateRecoverAccount(ctx, UpdateRecoverAccountParams{
			Uid:       uid,
			SecretKey: token,
		})
		if err != nil {
			return err
		}
		return q.UpdateUserActive(ctx, UpdateUserActiveParams{
			Uid: uid,
			IsActive: pgtype.Bool{
				Bool:  true,
				Valid: true,
			},
		})
	})
	return
}
