package app

import (
	"github.com/bannovdaniil/zeronews/internal/app/handler"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/list", handler.Home)
	app.Post("/edit/:Id", handler.EditNews)
}
