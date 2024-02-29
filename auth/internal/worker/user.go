package worker

import (
	"context"

	"github.com/hibiken/asynq"
)

const (
	TaskSendVerifyEmailUser    string = "task:send_verify_email"
	TaskSendRecoverAccountUser string = "task:send_recovery_account"
	TaskSendOTPUser            string = "task:send_otp_user"
)

type TaskUser interface {
	UserTaskSendVerificationEmail(ctx context.Context, payload *PayloadSendVerifyEmail, opts ...asynq.Option) error
	UserTaskSendRecoverAccount(ctx context.Context, payload *PayloadSendRecoverAccount, opts ...asynq.Option) error
}

type RedisTaskUser struct {
	client *asynq.Client
}

func NewRedisTaskDistributor(redisOpt asynq.RedisClientOpt) TaskUser {
	client := asynq.NewClient(redisOpt)
	return &RedisTaskUser{
		client: client,
	}
}
