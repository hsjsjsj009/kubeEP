package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/handler"
)

func buildRoute(handlers *handler.Handlers, router fiber.Router) {
	router.Use(
		cors.New(
			cors.Config{
				AllowHeaders: "Origin, Content-Type, Accept",
				AllowOrigins: "http://localhost:3000",
			},
		),
	)

	router.Route(
		"/gcp", func(router fiber.Router) {
			router.Route(
				"/register", func(router fiber.Router) {
					router.Post("/datacenter", handlers.GcpHandler.RegisterDatacenter)
					router.Post("/clusters", handlers.GcpHandler.RegisterClusterWithDatacenter)
				},
			)
			router.Get("/clusters", handlers.GcpHandler.GetClustersByDatacenterID)
		},
	)

	router.Route(
		"/clusters", func(router fiber.Router) {
			router.Get("/", handlers.ClusterHandler.GetAllRegisteredClusters)
			router.Get("/hpa", handlers.ClusterHandler.GetAllHPA)
		},
	)

	router.Route(
		"/event", func(router fiber.Router) {
			router.Post("/register", handlers.EventHandler.RegisterEvents)
		},
	)
}
