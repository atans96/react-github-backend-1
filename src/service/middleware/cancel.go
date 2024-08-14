package middleware

import (
	"context"
	"github.com/gofiber/fiber/v2"
)

func CancelReq(ctx *fiber.Ctx) error {
	c, cancel := context.WithCancel(context.Background())

	ctx.Locals("cancelFn", cancel)
	ctx.SetUserContext(c)

	err := ctx.Next()

	cancelFnWillBeCalled := ctx.Locals("cancelFn")

	if cancelFnWillBeCalled == nil {
		defer cancel()
	}

	return err
}
