package handler

import (
	"encoding/json"
	"github.com/bannovdaniil/zeronews/internal/app/model"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gopkg.in/reform.v1"
)

func GetList(context *fiber.Ctx, db *reform.DB, log *logrus.Logger) error {
	list, err := db.SelectAllFrom(model.NewsTable, "", nil)
	if err != nil {
		context.Status(fiber.StatusNoContent)
		log.Errorf("[Controller:getList] couldn't find any note %v", err)
		return nil
	}

	resp, err := json.Marshal(list)
	if err != nil {
		context.Status(fiber.StatusInternalServerError)
		log.Info("[Controller:getList] failed to marshal response body: %v", err)
		return nil
	}
	context.Set("Content-Type", "application/json")
	_, err = context.Write(resp)

	return context.SendString("Hello Root")
}
