package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/handler"
)

func buildRoute(handlers *handler.Handlers, router fiber.Router) {
	router.Route("/gcp", func(router fiber.Router) {
		router.Route("/register", func(router fiber.Router) {
			router.Post("/datacenter", handlers.GcpHandler.RegisterDatacenter)
			router.Post("/clusters", handlers.GcpHandler.RegisterClusterWithDatacenter)
		})
		router.Get("/clusters", handlers.GcpHandler.GetClustersByDatacenterID)
	})

	router.Route("/clusters", func(router fiber.Router) {
		router.Get("/", handlers.ClusterHandler.GetAllRegisteredClusters)
		router.Get("/hpa", handlers.ClusterHandler.GetAllHPA)
	})
}
