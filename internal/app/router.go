package app

import (
	"github.com/bannovdaniil/zeronews/internal/app/handler"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gopkg.in/reform.v1"
)

type NewsHandlers struct {
	db     *reform.DB
	logger *logrus.Logger
}

func (listHandler *NewsHandlers) listHandler(c *fiber.Ctx) error {
	return handler.GetList(c, listHandler.db, listHandler.logger)
}

func SetupRoutes(app *fiber.App, db *reform.DB, logger *logrus.Logger) {

	getListHandler := &NewsHandlers{db, logger}

	app.Get("/list", getListHandler.listHandler)
	app.Post("/edit/:Id", handler.EditNews)
}
