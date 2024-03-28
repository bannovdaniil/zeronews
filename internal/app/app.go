package app

import (
	"errors"
	config "github.com/bannovdaniil/zeronews/internal/app/config"
	"github.com/bannovdaniil/zeronews/internal/app/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func RunApp() {
	log := logrus.New()
	log.Info("Start server...")

	applicationConfig := config.LoadConfig(log)

	log.Info(applicationConfig)

	reformDB, err := repository.ConnectToDB(applicationConfig.DbUrl, log)

	if err != nil {
		log.Fatal(err)
	}
	server := fiber.New(fiber.Config{ErrorHandler: func(ctx *fiber.Ctx, err error) error {
		code := ctx.Response().StatusCode()
		var e *fiber.Error
		if errors.As(err, &e) {
			code = e.Code
		}

		err = ctx.Status(code).SendString(err.Error())
		if err != nil {
			return ctx.Status(e.Code).SendString(e.Message)
		}

		return nil
	},
	})
	SetupRoutes(server, reformDB, log)

	if err := server.Listen(":" + applicationConfig.Port); err != nil {
		log.Fatalf("Error: %s", err.Error())
		return
	}

}
