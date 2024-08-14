package routes

import "github.com/gofiber/fiber/v2"

func heavyTask() {
	for {
	}
}
func Test(c *fiber.Ctx) error {
	go heavyTask()
	return c.SendString("hi")
}
