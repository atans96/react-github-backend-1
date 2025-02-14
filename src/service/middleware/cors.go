package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"server_go/src/service"
	"strconv"
	"strings"
)

type Config struct {
	// Next defines a function to skip this middleware when returned true.
	//
	// Optional. Default: nil
	Next func(c *fiber.Ctx) bool

	// AllowOrigin defines a list of origins that may access the resource.
	//
	// Optional. Default value "*"
	AllowOrigins string

	// AllowMethods defines a list methods allowed when accessing the resource.
	// This is used in response to a preflight request.
	//
	// Optional. Default value "GET,POST,HEAD,PUT,DELETE,PATCH"
	AllowMethods string

	// AllowHeaders defines a list of request headers that can be used when
	// making the actual request. This is in response to a preflight request.
	//
	// Optional. Default value "".
	AllowHeaders string

	// AllowCredentials indicates whether or not the response to the request
	// can be exposed when the credentials flag is true. When used as part of
	// a response to a preflight request, this indicates whether or not the
	// actual request can be made using credentials.
	//
	// Optional. Default value false.
	AllowCredentials bool

	// ExposeHeaders defines a whitelist headers that clients are allowed to
	// access.
	//
	// Optional. Default value "".
	ExposeHeaders string

	// MaxAge indicates how long (in seconds) the results of a preflight request
	// can be cached.
	//
	// Optional. Default value 0.
	MaxAge int
}

// ConfigDefault is the default config
var defaultConfig = cors.Config{
	Next:         nil,
	AllowOrigins: "*",
	AllowMethods: strings.Join([]string{
		fiber.MethodGet,
		fiber.MethodPost,
		fiber.MethodHead,
		fiber.MethodPut,
		fiber.MethodDelete,
		fiber.MethodPatch,
	}, ","),
	AllowHeaders:     "",
	AllowCredentials: false,
	ExposeHeaders:    "",
	MaxAge:           0,
}

func CORSMiddleware(config ...cors.Config) fiber.Handler {
	// Set default config
	cfg := defaultConfig

	// Override config if provided
	if len(config) > 0 {
		cfg = config[0]

		// Set default values
		if cfg.AllowMethods == "" {
			cfg.AllowMethods = defaultConfig.AllowMethods
		}
		if cfg.AllowOrigins == "" {
			cfg.AllowOrigins = defaultConfig.AllowOrigins
		}
	}

	// Convert string to slice
	allowOrigins := strings.Split(strings.Replace(cfg.AllowOrigins, " ", "", -1), ",")

	// Strip white spaces
	allowMethods := strings.Replace(cfg.AllowMethods, " ", "", -1)
	allowHeaders := strings.Replace(cfg.AllowHeaders, " ", "", -1)
	exposeHeaders := strings.Replace(cfg.ExposeHeaders, " ", "", -1)

	// Convert int to string
	maxAge := strconv.Itoa(cfg.MaxAge)

	// Return new handler
	return func(c *fiber.Ctx) error {
		// Don't execute middleware if Next returns true
		if cfg.Next != nil && cfg.Next(c) {
			return c.Next()
		}

		// Get origin header
		origin := c.Get(fiber.HeaderOrigin)
		if len(origin) == 0 {
			return c.Status(fiber.StatusForbidden).SendString("CORS Error")
		}
		allowOrigin := ""

		// Check allowed origins
		for _, o := range allowOrigins {
			if o == "*" && cfg.AllowCredentials {
				allowOrigin = origin
				break
			}
			if o == "*" || o == origin {
				allowOrigin = o
				break
			}
			if strings.Compare(origin, o) == 0 {
				allowOrigin = origin
				break
			}
		}
		if len(allowOrigin) == 0 {
			return c.Status(fiber.StatusForbidden).SendString("CORS Error")
		}
		// Simple request
		req := c.Request()
		if service.Contains(strings.Split(allowMethods, ","), string(req.Header.Method())) {
			c.Vary(fiber.HeaderOrigin)
			c.Set(fiber.HeaderAccessControlAllowOrigin, allowOrigin)

			if cfg.AllowCredentials {
				c.Set(fiber.HeaderAccessControlAllowCredentials, "true")
			}
			if exposeHeaders != "" {
				c.Set(fiber.HeaderAccessControlExposeHeaders, exposeHeaders)
			}
			//return c.Next()
		} else {
			return c.Status(fiber.StatusForbidden).SendString("CORS Error")
		}

		// Preflight request
		c.Vary(fiber.HeaderOrigin)
		c.Vary(fiber.HeaderAccessControlRequestMethod)
		c.Vary(fiber.HeaderAccessControlRequestHeaders)
		c.Set(fiber.HeaderAccessControlAllowOrigin, allowOrigin)
		c.Set(fiber.HeaderAccessControlAllowMethods, allowMethods)

		// Set Allow-Credentials if set to true
		if cfg.AllowCredentials {
			c.Set(fiber.HeaderAccessControlAllowCredentials, "true")
		}

		// Set Allow-Headers if not empty
		if allowHeaders != "" {
			c.Set(fiber.HeaderAccessControlAllowHeaders, allowHeaders)
		} else {
			h := c.Get(fiber.HeaderAccessControlRequestHeaders)
			if h != "" {
				c.Set(fiber.HeaderAccessControlAllowHeaders, h)
			}
		}

		// Set MaxAge is set
		if cfg.MaxAge > 0 {
			c.Set(fiber.HeaderAccessControlMaxAge, maxAge)
		}

		// Send 204 No Content
		//return c.SendStatus(fiber.StatusNoContent)
		return c.Next()
	}
}
