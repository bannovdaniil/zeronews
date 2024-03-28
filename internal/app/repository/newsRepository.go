package repository

import (
	"github.com/bannovdaniil/zeronews/internal/app/model"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/sirupsen/logrus"
	"gopkg.in/reform.v1"
)

func FindAll(db *reform.DB, log *logrus.Logger) ([]model.News, error) {
	var args []interface{}
	newsTableList, err := db.SelectAllFrom(model.NewsTable, "", args...)
	if err != nil {
		log.Errorf("[Storage:findAll] %s", err)
		return nil, err
	}

	newsList := make([]model.News, 0, len(newsTableList))
	for _, noteStruct := range newsTableList {
		news := *noteStruct.(*model.News)

		categories, err := db.SelectAllFrom(model.CategoryTable, "WHERE news_id = $1", news.Id)
		if err != nil {
			log.Errorf("[Storage:findAll] Can't get categories %s", err)
			return nil, err
		}

		news.Categories = []int64{}
		for _, category := range categories {
			news.Categories = append(news.Categories, category.(*model.Category).CategoryId)
		}
		newsList = append(newsList, news)
	}

	return newsList, nil
}

func UpdateNews(db *reform.DB, log *logrus.Logger, news *model.News) (model.News, error) {
	newsRaw, err := db.SelectOneFrom(model.NewsTable, "WHERE id = $1", news.Id)
	if err != nil {
		log.Errorf("[Storage:UpdateNews]: %s", err)
		return model.News{}, err
	}

	newsEntity := newsRaw.(*model.News)
	newsEntity.Title = news.Title
	newsEntity.Content = news.Content
	newsEntity.Categories = news.Categories

	if err = db.Update(newsEntity); err != nil {
		log.Errorf("[Storage:UpdateNews]: %s", err)
		return model.News{}, err
	}

	return *newsEntity, nil
}
