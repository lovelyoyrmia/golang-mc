package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store interface {
	Querier
	ExecTx(ctx context.Context, fn func(*Queries) error) error
	CreateUserTx(ctx context.Context, req UserTxParams) error
	ValidateEmailTx(ctx context.Context, uid string, token string) error
	ValidateRecoverAccountTx(ctx context.Context, uid string, token string) error
}

type SQLStore struct {
	*Queries
	ConnPool *pgxpool.Pool
}

func NewStore(connPool *pgxpool.Pool) Store {
	return &SQLStore{
		Queries:  New(connPool),
		ConnPool: connPool,
	}
}

type UserTxParams struct {
	UserParams  CreateUserParams
	AfterCreate func(user User) error
}

func (store *SQLStore) ExecTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.ConnPool.Begin(ctx)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx error: %d", rbErr)
		}
		return err
	}

	return tx.Commit(ctx)
}
