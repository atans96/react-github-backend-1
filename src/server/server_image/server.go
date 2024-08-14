package server

import (
	"crypto/tls"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"os"
	"server_go/src/server/server_image/routes"
	"server_go/src/service/middleware"
)

func Run() {
	app := fiber.New(fiber.Config{IdleTimeout: 3600})
	trackMiddleware := middleware.RequestTrackerMiddleware{}
	trackMiddleware.Config.RequestID = "tracker_id"
	app.Use(trackMiddleware.New())
	app.Use(recover.New(recover.Config{EnableStackTrace: true}))
	app.Use(requestid.New(requestid.Config{ContextKey: "trackRequest"}))
	app.Use(middleware.CORSMiddleware(cors.Config{
		AllowOrigins:     "https://" + os.Getenv("CLIENT_HOST") + ":" + os.Getenv("CLIENT_PORT"),
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET,POST,OPTIONS",
		MaxAge:           -1,
		AllowCredentials: true,
	}))
	app.Use("/images_from_markdown", middleware.CancelReq)
	app.Options("/images_from_markdown", func(c *fiber.Ctx) error {
		return c.SendString("ok") // => https
	})
	app.Post("/images_from_markdown", routes.ImagesFromMarkdown)

	cert, err := tls.X509KeyPair([]byte(os.Getenv("localhost.crt")), []byte(os.Getenv("localhost.decrypted.key")))
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	// Create custom listener
	ln, err := tls.Listen("tcp", os.Getenv("GOLANG_HOST")+":"+os.Getenv("GOLANG_PORT_IMG"), tlsConfig)
	if err != nil {
		panic(err)
	}
	fmt.Println(app.Listener(ln))
}
