package user

import (
	"github.com/Foedie/foedie-server-v2/gateway/internal/config"
	"github.com/Foedie/foedie-server-v2/gateway/pkg/auth"
	"github.com/Foedie/foedie-server-v2/gateway/pkg/user/routes"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(route chi.Router, c *config.Config, authMiddleware *auth.AuthMiddlewareConfig) *ServiceClient {
	svc := &ServiceClient{
		Client: InitServiceClient(c),
	}

	userRoutes := routes.NewUserRoutes(svc.Client)
	route.Route("/user", func(r chi.Router) {
		r.Use(authMiddleware.AuthMiddleware)
		r.Get("/", userRoutes.GetUser)
		r.Put("/", userRoutes.UpdateUser)
		r.Delete("/", userRoutes.DeleteUser)
	})

	return svc
}
