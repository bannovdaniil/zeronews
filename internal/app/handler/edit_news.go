package handler

import (
	"github.com/gofiber/fiber/v2"
)

func EditNews(ctx *fiber.Ctx) error {
	newsId := ctx.Params("Id", "key was not provided")
	//news := new(model.News)
	//if err := ctx.BodyParser(news); err != nil {
	//	return ctx.Status(fiber.StatusInternalServerError).JSON(
	//		fiber.Map{"message": err.Error()})
	//}
	//
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"get id": newsId})
}
