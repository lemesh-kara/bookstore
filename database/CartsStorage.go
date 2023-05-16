package database

import (
	"bookstore/model"
	"errors"
	"strconv"

	"github.com/jinzhu/gorm"
)

type CartsStorage struct {
	db *gorm.DB
}

func NewCartsStorage(db *gorm.DB) *CartsStorage {
	return &CartsStorage{
		db: db,
	}
}

func (storage *CartsStorage) AddCart(cart *model.Cart) error {
	err := storage.db.Create(cart).Error
	if err != nil {
		return errors.New("CartsStorage error: " + err.Error())
	}
	return nil
}

func (storage *CartsStorage) GetCart(id uint64) (model.Cart, error) {
	var cart model.Cart
	err := storage.db.Preload("User").Preload("Book").First(&cart, "id = ?", id).Error

	if err != nil {
		return model.Cart{}, errors.New("CartsStorage error: " + err.Error())
	}

	resetUserSecrets(&cart.User)

	return cart, nil
}

func (storage *CartsStorage) DeleteCart(id uint64) error {
	cart, err := storage.GetCart(id)
	if err != nil {
		return errors.New("CartsStorage error: " + err.Error())
	}
	err = storage.db.Delete(&cart).Error
	if err != nil {
		return errors.New("CartsStorage error: " + err.Error())
	}
	return nil
}

func (storage *CartsStorage) GetAllCarts(limit, offset uint64) ([]model.Cart, error) {
	var cartsList []model.Cart
	var err error

	if limit > 0 || offset > 0 {
		err = storage.db.Offset(offset).Limit(limit).Find(&cartsList).Error
	} else {
		err = storage.db.Find(&cartsList).Error
	}

	if err != nil {
		return []model.Cart{}, errors.New("CartsStorage error: " + err.Error())
	}
	return cartsList, nil
}

func (storage *CartsStorage) GetCartsByUserId(userID uint64, limit, offset uint64) ([]model.Cart, error) {
	var carts []model.Cart
	var err error

	if limit > 0 || offset > 0 {
		err = storage.db.Debug().Offset(offset).Limit(limit).Preload("User").Preload("Book").Where("user_id = ?", userID).Find(&carts).Error
	} else {
		err = storage.db.Debug().Preload("User").Preload("Book").Where("user_id = ?", userID).Find(&carts).Error
	}

	if err != nil {
		return []model.Cart{}, errors.New("CartsStorage error: " + err.Error() + " user id " + strconv.FormatUint(userID, 10))
	}

	for index, _ := range carts {
		resetUserSecrets(&carts[index].User)
	}

	return carts, nil
}

func (storage *CartsStorage) GetCartsByBookId(bookID uint64, limit, offset uint64) ([]model.Cart, error) {
	var carts []model.Cart
	var err error

	if limit > 0 || offset > 0 {
		err = storage.db.Offset(offset).Limit(limit).Preload("User").Preload("Book").Where("book_id = ?", bookID).Find(&carts).Error
	} else {
		err = storage.db.Preload("User").Preload("Book").Where("book_id = ?", bookID).Find(&carts).Error
	}

	if err != nil {
		return []model.Cart{}, errors.New("CartsStorage error: " + err.Error())
	}

	for index, _ := range carts {
		resetUserSecrets(&carts[index].User)
	}

	return carts, nil
}
