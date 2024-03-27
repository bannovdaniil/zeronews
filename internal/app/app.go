package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func RunApp() {
	logrus.Info("Start server...")
	server := fiber.New()

	SetupRoutes(server)

	if err := server.Listen(":3000"); err != nil {
		logrus.Fatalf("Error: %s", err.Error())
		return
	}

}
