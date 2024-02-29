package main

import (
	"context"
	"errors"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/Foedie/foedie-server-v2/auth/pkg/logger"
	"github.com/Foedie/foedie-server-v2/email/domain/pb"
	"github.com/Foedie/foedie-server-v2/email/domain/services"
	"github.com/Foedie/foedie-server-v2/email/pkg/config"
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

	ctx, stop := signal.NotifyContext(context.Background(), interruptSignal...)
	defer stop()

	waitGroup, ctx := errgroup.WithContext(ctx)

	runGrpcServer(waitGroup, ctx, c)

	err = waitGroup.Wait()
	if err != nil {
		log.Fatal().Err(err).Msg("error from wait group")
	}
}

func runGrpcServer(
	waitGroup *errgroup.Group,
	ctx context.Context,
	c config.Config,
) {
	grpcLogger := grpc.UnaryInterceptor(logger.GrpcLogger)

	server := services.NewServer(c)

	lis, err := net.Listen("tcp", c.Port)

	if err != nil {
		log.Fatal().Msgf("Failed to listen: %d", err)
	}

	grpcServer := grpc.NewServer(grpcLogger)
	pb.RegisterEmailServiceServer(grpcServer, server)

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
