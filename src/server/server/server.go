package server

import (
	"github.com/gofiber/fiber/v2"
	"os"
	"server_go/src/server/server/routes"
	"server_go/src/service"
	"server_go/src/service/middleware"
)

func Run() {
	app := fiber.New(fiber.Config{IdleTimeout: 3600})
	x := &service.Server{App: app}
	middlewareService := middleware.MiddlewareService{Server: x}
	middlewareService.GlobalMiddleware()
	middlewareService.RouteMiddleware()
	s := &routes.Routes{Server: x}
	s.StartRoutes()
	x.StartServer(os.Getenv("GOLANG_HOST") + ":" + os.Getenv("GOLANG_PORT"))
}
