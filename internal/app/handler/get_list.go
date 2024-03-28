package handler

import (
	"encoding/json"
	"github.com/bannovdaniil/zeronews/internal/app/model"
	"github.com/bannovdaniil/zeronews/internal/app/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gopkg.in/reform.v1"
)

type responseList struct {
	Success bool         `json:"Success"`
	News    []model.News `json:"News"`
}

func GetList(context *fiber.Ctx, db *reform.DB, log *logrus.Logger) error {
	errorResponse := responseList{
		false,
		make([]model.News, 0),
	}

	newsList, err := repository.FindAll(db, log)
	if err != nil {
		log.Errorf("[Controller:getList] %s", err)
		_ = sendJson(errorResponse, context, fiber.StatusInternalServerError)
		return err
	}

	response := responseList{
		true,
		newsList,
	}

	err = sendJson(response, context, fiber.StatusOK)

	return err
}

func sendJson(response responseList, context *fiber.Ctx, serverStatus int) error {
	resp, err := json.Marshal(response)
	if err != nil {
		context.Status(fiber.StatusInternalServerError)
		return err
	}

	context.Status(serverStatus)
	context.Set("Content-Type", "application/json")
	_, err = context.Write(resp)

	return err
}
