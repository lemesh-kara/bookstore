package database

import (
	"bookstore/model"
	"errors"

	"github.com/jinzhu/gorm"
)

type NewsStorage struct {
	db *gorm.DB
}

func NewNewsStorage(db *gorm.DB) *NewsStorage {
	return &NewsStorage{
		db: db,
	}
}

func (storage *NewsStorage) GetNews(id uint64) (model.News, error) {
	var aaaa model.News

	err := storage.db.First(&aaaa, "id = ?", id).Error

	if err != nil {
		return model.News{}, errors.New("NewsStorage error: " + err.Error())
	}

	return aaaa, nil
}

func (storage *NewsStorage) GetAllNews() ([]model.News, error) {
	var news []model.News
	var err error

	err = storage.db.Find(&news).Error

	if err != nil {
		return []model.News{}, errors.New("NewsStorage error: " + err.Error())
	}
	return news, nil
}
