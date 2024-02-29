package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Foedie/foedie-server-v2/gateway/internal/config"
	"github.com/Foedie/foedie-server-v2/gateway/internal/logger"
	"github.com/Foedie/foedie-server-v2/gateway/pkg/auth"
	"github.com/Foedie/foedie-server-v2/gateway/pkg/user"
	"github.com/go-chi/chi/v5"
	chiMid "github.com/go-chi/chi/v5/middleware"
	"github.com/go-openapi/runtime/middleware"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

var interruptSignal = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

func main() {

	c, err := config.LoadConfig()

	if err != nil {
		log.Fatal().Msgf("cannot loag config: %d", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), interruptSignal...)
	defer stop()

	waitGroup, ctx := errgroup.WithContext(ctx)

	runGatewayServer(waitGroup, ctx, c)

	err = waitGroup.Wait()
	if err != nil {
		log.Fatal().Err(err).Msg("error from wait group")
	}
}

func runGatewayServer(waitGroup *errgroup.Group, ctx context.Context, c config.Config) {
	route := chi.NewRouter()

	route.Use(logger.HTTPLogger)
	route.Use(chiMid.Recoverer)

	route.Route("/api", func(v1 chi.Router) {
		v1.Route("/v1", func(r chi.Router) {
			authSvc := auth.RegisterRoutes(r, &c)

			authMiddleware := auth.InitAuthMiddleware(authSvc)

			user.RegisterRoutes(r, &c, &authMiddleware)
		})

		// DOCS

		ops := middleware.RedocOpts{SpecURL: "/foedie.swagger.json"}
		sh := middleware.Redoc(ops, nil)
		route.Method(http.MethodGet, "/docs", sh)
		fs := http.FileServer(http.Dir("./docs/swagger"))
		route.Method(http.MethodGet, "/foedie.swagger.json", fs)

	})

	s := &http.Server{
		Addr:         c.Port,
		Handler:      http.TimeoutHandler(route, time.Duration(time.Second*30), "server timeout"),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  5 * time.Second,
	}

	waitGroup.Go(func() error {
		log.Info().Msgf("Starting HTTP Gateway on port %s", s.Addr)
		err := s.ListenAndServe()
		if err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return nil
			}
			log.Fatal().Msgf("Error starting server : %s\n", err)
			return err
		}
		return nil
	})

	waitGroup.Go(func() error {
		<-ctx.Done()
		log.Info().Msg("graceful shutdown HTTP gateway server")

		err := s.Shutdown(context.Background())
		if err != nil {
			log.Error().Msg("failed to shutdown http gateway server")
			return err
		}
		log.Info().Msg("HTTP Gateway server is stopped")
		return nil
	})
}
