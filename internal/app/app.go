package app

import (
	config "github.com/bannovdaniil/zeronews/internal/app/config"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func RunApp() {
	log := logrus.New()
	log.Info("Start server...")

	applicationConfig := config.LoadConfig(log)

	log.Println(applicationConfig)

	server := fiber.New()
	SetupRoutes(server)

	addr := applicationConfig.Host + ":" + applicationConfig.Port

	if err := server.Listen(addr); err != nil {
		log.Fatalf("Error: %s", err.Error())
		return
	}

}
