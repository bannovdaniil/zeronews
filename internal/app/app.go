package app

import (
	config "github.com/bannovdaniil/zeronews/internal/app/config"
	"github.com/bannovdaniil/zeronews/internal/app/database"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func RunApp() {
	log := logrus.New()
	log.Info("Start server...")

	applicationConfig := config.LoadConfig(log)

	log.Info(applicationConfig)

	reformDB, err := database.ConnectToDB(applicationConfig.DbUrl, log)

	if err != nil {
		log.Fatal(err)
	}
	server := fiber.New()
	SetupRoutes(server, reformDB, log)

	if err := server.Listen(":" + applicationConfig.Port); err != nil {
		log.Fatalf("Error: %s", err.Error())
		return
	}

}
