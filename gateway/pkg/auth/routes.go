package auth

import (
	"github.com/Foedie/foedie-server-v2/gateway/internal/config"
	"github.com/Foedie/foedie-server-v2/gateway/pkg/auth/routes"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, c *config.Config) *ServiceClient {
	svc := &ServiceClient{
		Client: InitServiceClient(c),
	}

	authRoutes := routes.NewAuthRoutes(svc.Client)
	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", authRoutes.Login)
		r.Post("/register", authRoutes.Register)
		r.Post("/verify-email", authRoutes.VerifyEmail)
		r.Get("/verify-email/{token}", authRoutes.ValidateEmail)
		r.Post("/recover-account", authRoutes.RecoverAccount)
		r.Get("/recover-account/{token}", authRoutes.ValidateRecoverAccount)
	})

	return svc
}
