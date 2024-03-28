package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	model "github.com/bannovdaniil/zeronews/internal/app/model"
	"github.com/bannovdaniil/zeronews/internal/app/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gopkg.in/reform.v1"
	"strconv"
)

func UpdateNews(context *fiber.Ctx, db *reform.DB, log *logrus.Logger) error {
	rawNewsId := context.Params("Id")

	newsId, err := strconv.ParseInt(rawNewsId, 10, 64)
	if err != nil || newsId <= 0 {
		context.Status(fiber.StatusBadRequest)
		log.Info("[Controller:UpdateNews] invalid id: %s", err)
		return err
	}

	var updateNews *model.News
	if err = json.Unmarshal(context.Body(), &updateNews); err != nil {
		context.Status(fiber.StatusBadRequest)
		log.Info("[Controller:UpdateNews] bad unmarshal request body")
		return err
	}
	if newsId != updateNews.Id {
		context.Status(fiber.StatusBadRequest)
		log.Info("[Controller:UpdateNews] bad request, ID is wrong")
		return nil
	}

	news, err := repository.UpdateNews(db, log, updateNews)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			context.Status(fiber.StatusNotFound)
		} else {
			context.Status(fiber.StatusBadRequest)
		}
		log.Info("[Controller:UpdateNews] database error: %s", err)
		return err
	}

	resp, err := json.Marshal(news)
	context.Status(fiber.StatusOK)
	context.Set("Content-Type", "application/json")
	_, err = context.Write(resp)

	return err
}
