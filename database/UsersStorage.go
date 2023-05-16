package database

import (
	"bookstore/model"
	"errors"

	"github.com/jinzhu/gorm"
)

type UsersStorage struct {
	db *gorm.DB
}

func NewUsersStorage(db *gorm.DB) *UsersStorage {
	return &UsersStorage{
		db: db,
	}
}

func (storage *UsersStorage) AddUser(user *model.User) error {
	err := storage.db.Create(user).Error
	if err != nil {
		return errors.New("UsersStorage error: " + err.Error())
	}
	return nil
}

func (storage *UsersStorage) GetUser(id uint64) (model.User, error) {
	var user model.User
	err := storage.db.First(&user, "id = ?", id).Error

	if err != nil {
		return model.User{}, errors.New("UsersStorage error: " + err.Error())
	}

	resetUserPassword(&user)

	return user, nil
}

func (storage *UsersStorage) DeleteUser(id uint64) error {
	user, err := storage.GetUser(id)
	if err != nil {
		return errors.New("UsersStorage error: " + err.Error())
	}
	err = storage.db.Delete(&user).Error
	if err != nil {
		return errors.New("UsersStorage error: " + err.Error())
	}
	return nil
}

func (storage *UsersStorage) GetAllUsers(limit, offset uint64) ([]model.User, error) {
	var usersList []model.User
	var err error

	if limit > 0 || offset > 0 {
		err = storage.db.Offset(offset).Limit(limit).Find(&usersList).Error
	} else {
		err = storage.db.Find(&usersList).Error
	}

	if err != nil {
		return []model.User{}, errors.New("UsersStorage error: " + err.Error())
	}

	for index, _ := range usersList {
		resetUserPassword(&usersList[index])
	}

	return usersList, nil
}

func (storage *UsersStorage) GetUserByUsername(username string) (model.User, error) {
	var user model.User
	err := storage.db.First(&user, "username = ?", username).Error

	if err != nil {
		return model.User{}, errors.New("UsersStorage error: " + err.Error())
	}

	return user, nil
}

func (storage *UsersStorage) GetFullUser(id uint64) (model.User, error) {
	var user model.User
	err := storage.db.Preload("Reviews").First(&user, "id = ?", id).Error

	if err != nil {
		return model.User{}, errors.New("UsersStorage error: " + err.Error())
	}

	resetUserPassword(&user)

	return user, nil
}
