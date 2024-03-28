package app

import (
	config "github.com/bannovdaniil/zeronews/internal/app/config"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func RunApp() {
	logrus.Info("Start server...")

	applicationConfig := config.LoadConfig()

	logrus.Println(applicationConfig)

	server := fiber.New()
	SetupRoutes(server)

	if err := server.Listen(":" + applicationConfig.Port); err != nil {
		logrus.Fatalf("Error: %s", err.Error())
		return
	}

}
