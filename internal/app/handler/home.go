package handler

import "github.com/gofiber/fiber/v2"

func Home(context *fiber.Ctx) error {
	return context.SendString("Hello Root")
}
