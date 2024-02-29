package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

func (store *SQLStore) UpdateUserTx(ctx context.Context, user GetUserRow, req UpdateUserParams) error {
	return store.ExecTx(ctx, func(q *Queries) error {
		user1, err := q.UpdateUser(ctx, req)
		if err != nil {
			return err
		}

		if user.Email != user1.Email {
			return q.UpdateUserVerified(ctx, UpdateUserVerifiedParams{
				Uid: user1.Uid,
				IsVerified: pgtype.Bool{
					Bool:  false,
					Valid: true,
				},
			})
		}
		return nil
	})
}
