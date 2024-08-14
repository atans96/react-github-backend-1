package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/websocket/v2"
	"os"
	"server_go/src/service"
)

type MiddlewareService struct {
	*service.Server
}

func (app *MiddlewareService) GlobalMiddleware() {
	trackMiddleware := RequestTrackerMiddleware{}
	trackMiddleware.Config.RequestID = "tracker_id"
	app.Use(trackMiddleware.New())
	app.Use(recover.New(recover.Config{EnableStackTrace: true}))
	app.Use(requestid.New(requestid.Config{ContextKey: "trackRequest"}))
	app.Use(CORSMiddleware(cors.Config{
		AllowOrigins:     "https://" + os.Getenv("CLIENT_HOST") + ":" + os.Getenv("CLIENT_PORT"),
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET,POST,OPTIONS",
		MaxAge:           -1,
		AllowCredentials: true,
	}))
	app.Use("/images_from_markdown", CancelReq)
	app.Use("/graphqlws", func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})
}
func (app *MiddlewareService) RouteMiddleware() {
}
