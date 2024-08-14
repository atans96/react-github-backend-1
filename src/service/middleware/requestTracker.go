package middleware

import (
	"fmt"
	"github.com/SkyAPM/go2sky"
	"github.com/gofiber/fiber/v2"
)

// Config defines the config for middleware.
type RequestTrackerConfig struct {
	// Next defines a function to skip this middleware when returned true.
	//
	// Optional. Default: nil
	Next      func(c *fiber.Ctx) bool
	RequestID string
}

// ConfigDefault is the default config
var ConfigDefault = RequestTrackerConfig{
	Next:      nil,
	RequestID: "",
}

// Helper function to set default values
func configDefault(config ...RequestTrackerConfig) RequestTrackerConfig {
	// Return default config if nothing provided
	if len(config) < 1 {
		return ConfigDefault
	}

	// Override default config
	cfg := config[0]
	return cfg
}

type RequestTrackerMiddleware struct {
	reporter go2sky.Reporter
	tracer   *go2sky.Tracer
	Config   RequestTrackerConfig
}

func (router *RequestTrackerMiddleware) New() fiber.Handler {
	// Set default config
	cfg := configDefault(router.Config)
	// Return new handler
	return func(c *fiber.Ctx) error {
		var err error
		rid := fmt.Sprintf("%v", c.Locals(cfg.RequestID))
		router.tracer, err = go2sky.NewTracer(rid, go2sky.WithReporter(router.reporter))
		if router.tracer == nil || err != nil {
			panic("requestTracker middleware error")
		}
		// Don't execute middleware if Next returns true
		if cfg.Next != nil && cfg.Next(c) {
			return c.Next()
		}
		return c.Next()
	}
}
