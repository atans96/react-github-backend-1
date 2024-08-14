package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"os"
	"server_go/src/service"
)

type Routes struct {
	*service.Server
}

func (app *Routes) StartRoutes() {
	// Routes
	app.Static("/", "files")
	app.Get("/test", Test)
	start := Starts{}
	app.Get("/start", start.Start)
	app.Options("/authenticate", func(c *fiber.Ctx) error {
		return c.SendString("ok") // => https
	})
	app.Post("/authenticate", Register)
	app.Get("/graphqlws", WS())
	app.Get("/rssFeed", RSSFeed)
	app.Options("/graphql/", func(c *fiber.Ctx) error {
		return c.SendString("ok") // => https
	})
	app.Post("/graphql/", proxy.Forward(os.Getenv("GRAPHQL_ADDRESS")+"/graphql/"))
	app.Options("/server_uwebsocket/end_of_session", func(c *fiber.Ctx) error {
		return c.SendString("ok") // => https
	})
	app.Post("/server_uwebsocket/end_of_session", proxy.Forward(os.Getenv("UWEBSOCKET_HOST")+":"+os.Getenv("UWEBSOCKET_PORT")+"/server_uwebsocket/end_of_session"))
	app.Options("/server_uwebsocket/auth_graphql", func(c *fiber.Ctx) error {
		return c.SendString("ok") // => https
	})
	app.Post("/server_uwebsocket/auth_graphql", proxy.Forward(os.Getenv("UWEBSOCKET_HOST")+":"+os.Getenv("UWEBSOCKET_PORT")+"/server_uwebsocket/auth_graphql"))
	app.Options("/server_uwebsocket/start_of_session", func(c *fiber.Ctx) error {
		return c.SendString("ok") // => https
	})
	app.Post("/server_uwebsocket/start_of_session", proxy.Forward(os.Getenv("UWEBSOCKET_HOST")+":"+os.Getenv("UWEBSOCKET_PORT")+"/server_uwebsocket/start_of_session"))
	app.Get("/server_uwebsocket/getTokenGQL", func(c *fiber.Ctx) error {
		url := os.Getenv("UWEBSOCKET_HOST") + ":" + os.Getenv("UWEBSOCKET_PORT") + "/server_uwebsocket/getTokenGQL?username=" + c.Query("username")
		if err := proxy.Do(c, url); err != nil {
			return err
		}
		// Remove Server header from response
		c.Response().Header.Del(fiber.HeaderServer)
		return nil
	})
	app.Options("/server_uwebsocket/setTokenGQL", func(c *fiber.Ctx) error {
		return c.SendString("ok") // => https
	})
	app.Post("/server_uwebsocket/setTokenGQL", proxy.Forward(os.Getenv("UWEBSOCKET_HOST")+":"+os.Getenv("UWEBSOCKET_PORT")+"/server_uwebsocket/setTokenGQL"))
	app.Options("/server_python/python_crawler", func(c *fiber.Ctx) error {
		return c.SendString("ok") // => https
	})
	app.Post("/server_python/python_crawler", proxy.Forward(os.Getenv("PYTHON_ADDRESS")+"/server_python/python_crawler"))
	app.Get("/server_uwebsocket/destroyTokenGQL", proxy.Forward(os.Getenv("UWEBSOCKET_HOST")+":"+os.Getenv("UWEBSOCKET_PORT")+"/server_uwebsocket/destroyTokenGQL"))
	app.Get("/server_uwebsocket/destroyToken", proxy.Forward(os.Getenv("UWEBSOCKET_HOST")+":"+os.Getenv("UWEBSOCKET_PORT")+"/server_uwebsocket/destroyToken"))
	app.Get("/server_uwebsocket/markdown", func(c *fiber.Ctx) error {
		url := os.Getenv("UWEBSOCKET_HOST") + ":" + os.Getenv("UWEBSOCKET_PORT") + "/server_uwebsocket/markdown?full_name=" + c.Query("full_name") + "&branch=" + c.Query("branch") + "&username=" + c.Query("username")
		if err := proxy.Do(c, url); err != nil {
			return err
		}
		// Remove Server header from response
		c.Response().Header.Del(fiber.HeaderServer)
		return nil
	})
}
