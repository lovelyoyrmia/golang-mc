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

type PayloadSendVerifyEmail struct {
	Email string `json:"email"`
}

// UserTaskSendVerificationEmail implements TaskUser.
func (user *RedisTaskUser) UserTaskSendVerificationEmail(ctx context.Context, payload *PayloadSendVerifyEmail, opts ...asynq.Option) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshall json: %w", err)
	}

	task := asynq.NewTask(TaskSendVerifyEmailUser, jsonPayload, opts...)
	info, err := user.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue user: %w", err)
	}

	log.Info().Int("max_retry", info.MaxRetry).Str("queue", info.Queue).Bytes("payload", info.Payload).Msg("Enqueue task...")
	return nil
}

// ProcessTaskSendVerificationEmail implements TaskProcessor.
func (processor *RedisTaskProcessor) ProcessTaskSendVerificationEmail(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendVerifyEmail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshall payload: %w", err)
	}

	user, err := processor.store.GetUserByEmailOrUsername(ctx, payload.Email)

	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	randmString := utils.RandomString(32)
	token := fmt.Sprintf("%s__%s", randmString, user.Uid)

	var verifyEmail db.VerifyEmail
	err = processor.store.ExecTx(ctx, func(q *db.Queries) error {
		verifyEmail, err = q.CreateVerifyEmail(ctx, db.CreateVerifyEmailParams{
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
		return fmt.Errorf("failed to create verify email: %w", err)
	}

	url := fmt.Sprintf("%s/verify-email/%s", processor.config.HOST, verifyEmail.SecretKey)

	res, err := processor.emailSvc.SendEmail(
		verifyEmail.Email,
		"Foedie Verification Email",
		url,
		verificationEmail,
	)

	if err != nil {
		return fmt.Errorf("failed to send verify email: %w", err)
	}

	log.Info().Msg(res.GetMessage())
	return nil
}
