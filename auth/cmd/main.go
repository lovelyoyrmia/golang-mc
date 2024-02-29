package main

import (
	"context"
	"errors"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/Foedie/foedie-server-v2/auth/domain/clients"
	"github.com/Foedie/foedie-server-v2/auth/domain/pb"
	"github.com/Foedie/foedie-server-v2/auth/domain/services"
	"github.com/Foedie/foedie-server-v2/auth/internal/db"
	"github.com/Foedie/foedie-server-v2/auth/internal/worker"
	"github.com/Foedie/foedie-server-v2/auth/pkg/config"
	"github.com/Foedie/foedie-server-v2/auth/pkg/logger"
	"github.com/Foedie/foedie-server-v2/auth/pkg/token"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var interruptSignal = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

func main() {
	c, err := config.LoadConfig()

	if err != nil {
		log.Fatal().Msgf("failed to  load config: %d", err)
	}

	redisOpt := asynq.RedisClientOpt{
		Addr: c.RedisAddr,
	}

	taskUser := worker.NewRedisTaskDistributor(redisOpt)

	ctx, stop := signal.NotifyContext(context.Background(), interruptSignal...)
	defer stop()

	database, err := db.NewDatabase(ctx, c)

	if err != nil {
		log.Fatal().Msgf("failed to connect database: %d", err)
	}

	store := db.NewStore(database.DB)

	tokenMaker, err := token.NewPasetoMaker(c)

	if err != nil {
		log.Fatal().Msgf("cannot create token %d", err)
	}

	waitGroup, ctx := errgroup.WithContext(ctx)

	runTaskProcessor(waitGroup, ctx, redisOpt, store, c)
	runGrpcServer(waitGroup, ctx, store, tokenMaker, c, taskUser)

	err = waitGroup.Wait()
	if err != nil {
		log.Fatal().Err(err).Msg("error from wait group")
	}
}

func runGrpcServer(
	waitGroup *errgroup.Group,
	ctx context.Context,
	store db.Store,
	tokenMaker token.Maker,
	c config.Config,
	taskUser worker.TaskUser,
) {
	grpcLogger := grpc.UnaryInterceptor(logger.GrpcLogger)
	server := services.NewServer(store, tokenMaker, taskUser)

	lis, err := net.Listen("tcp", c.Port)

	if err != nil {
		log.Fatal().Msgf("Failed to listen: %d", err)
	}

	grpcServer := grpc.NewServer(grpcLogger)
	pb.RegisterAuthServiceServer(grpcServer, server)

	reflection.Register(grpcServer)

	waitGroup.Go(func() error {
		log.Info().Msgf("Starting GRPC server on port: %s", lis.Addr().String())
		err = grpcServer.Serve(lis)
		if err != nil {
			if errors.Is(err, grpc.ErrServerStopped) {
				return nil
			}
			log.Fatal().Msgf("cannot start GRPC server: %d", err)
			return err
		}
		return nil
	})

	waitGroup.Go(func() error {
		<-ctx.Done()
		log.Info().Msg("graceful shutdown gRPC server")

		grpcServer.GracefulStop()
		log.Info().Msg("gRPC server is stopped")
		return nil
	})
}

func runTaskProcessor(
	waitGroup *errgroup.Group,
	ctx context.Context,
	redisOpt asynq.RedisClientOpt,
	store db.Store,
	config config.Config,
) {
	emailSvc := clients.InitEmailServiceClient(config.EmailSvcUrl)
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store, config, emailSvc)
	err := taskProcessor.Start()
	if err != nil {
		log.Fatal().Msgf("cannot start processor %d", err)
	}

	waitGroup.Go(func() error {
		<-ctx.Done()
		log.Info().Msg("graceful shutdown task processor")

		taskProcessor.Shutdown()
		log.Info().Msg("task processor is stopped")
		return nil
	})
}
