package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Foedie/foedie-server-v2/auth/internal/db"
	"github.com/Foedie/foedie-server-v2/auth/pkg/utils"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog/log"
)

type PayloadSendRecoverAccount struct {
	Email    string `json:"email"`
	Username string `json:"username"`
}

// UserTaskSendRecoverAccount implements TaskUser.
func (user *RedisTaskUser) UserTaskSendRecoverAccount(ctx context.Context, payload *PayloadSendRecoverAccount, opts ...asynq.Option) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshall json: %w", err)
	}

	task := asynq.NewTask(TaskSendRecoverAccountUser, jsonPayload, opts...)
	info, err := user.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue user: %w", err)
	}

	log.Info().Int("max_retry", info.MaxRetry).Str("queue", info.Queue).Bytes("payload", info.Payload).Msg("Enqueue task...")
	return nil
}

// ProcessTaskSendRecoverAccount implements TaskProcessor.
func (processor *RedisTaskProcessor) ProcessTaskSendRecoverAccount(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendRecoverAccount
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshall payload: %w", err)
	}

	user, err := processor.store.GetUserByEmailAndUsername(ctx, db.GetUserByEmailAndUsernameParams{
		Email:    payload.Email,
		Username: payload.Username,
	})

	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	randmString := utils.RandomString(32)
	token := fmt.Sprintf("%s__%s", randmString, user.Uid)

	var recoverAcc db.RecoverAccount
	err = processor.store.ExecTx(ctx, func(q *db.Queries) error {
		recoverAcc, err = q.CreateRecoveryAccount(ctx, db.CreateRecoveryAccountParams{
			ID:        uuid.New().String(),
			Uid:       user.Uid,
			SecretKey: token,
			Email:     payload.Email,
			ExpiredAt: pgtype.Timestamp{
				Time:  time.Now().Add(time.Minute * 15),
				Valid: true,
			},
		})
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to create recovery account: %w", err)
	}

	url := fmt.Sprintf("%s/recover-account/%s", "", recoverAcc.SecretKey)

	res, err := processor.emailSvc.SendEmail(
		recoverAcc.Email,
		"Foedie Recover Account",
		url,
		recoverAccount,
	)

	if err != nil {
		return fmt.Errorf("failed to send verify email: %w", err)
	}

	log.Info().Msg(res.GetMessage())
	return nil
}
