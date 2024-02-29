package worker

import (
	"context"

	"github.com/Foedie/foedie-server-v2/auth/domain/clients"
	"github.com/Foedie/foedie-server-v2/auth/internal/db"
	"github.com/Foedie/foedie-server-v2/auth/pkg/config"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const (
	QueueCritical string = "critical"
	QueueDefault  string = "default"
)

const (
	recoverAccount    = "recover-account"
	verificationEmail = "verify-email"
)

type TaskProcessor interface {
	Start() error
	Shutdown()
	ProcessTaskSendVerificationEmail(ctx context.Context, task *asynq.Task) error
	ProcessTaskSendRecoverAccount(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	server   *asynq.Server
	store    db.Store
	config   config.Config
	emailSvc clients.EmailServiceClient
}

func NewRedisTaskProcessor(
	redisOpt asynq.RedisClientOpt,
	store db.Store,
	config config.Config,
	emailSvc clients.EmailServiceClient,
) TaskProcessor {
	server := asynq.NewServer(redisOpt, asynq.Config{
		Queues: map[string]int{
			QueueCritical: 10,
			QueueDefault:  5,
		},
		ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
			log.Error().Err(err).Str("type", task.Type()).Bytes("payload", task.Payload()).Msg("process task failed")
		}),
	})

	return &RedisTaskProcessor{
		server:   server,
		store:    store,
		config:   config,
		emailSvc: emailSvc,
	}
}

// Start implements TaskProcessor.
func (processor *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()

	mux.HandleFunc(TaskSendVerifyEmailUser, processor.ProcessTaskSendVerificationEmail)
	mux.HandleFunc(TaskSendRecoverAccountUser, processor.ProcessTaskSendRecoverAccount)
	return processor.server.Start(mux)
}

// Shutdown implements TaskProcessor.
func (processor *RedisTaskProcessor) Shutdown() {
	processor.server.Shutdown()
}
