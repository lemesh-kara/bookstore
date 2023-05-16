package database

import (
	"errors"

	"bookstore/model"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type BooksStorage struct {
	db *gorm.DB
}

func NewBooksStorage(db *gorm.DB) *BooksStorage {
	return &BooksStorage{
		db: db,
	}
}

func (storage *BooksStorage) AddBook(book *model.Book) error {
	err := storage.db.Create(book).Error
	if err != nil {
		return errors.New("BooksStorage error: " + err.Error())
	}
	return nil
}

func (storage *BooksStorage) GetBook(id uint64) (model.Book, error) {
	var book model.Book
	err := storage.db.First(&book, "id = ?", id).Error

	if err != nil {
		return model.Book{}, errors.New("BooksStorage error: " + err.Error())
	}
	return book, nil
}

func (storage *BooksStorage) DeleteBook(id uint64) error {
	book, err := storage.GetBook(id)
	if err != nil {
		return err
	}
	err = storage.db.Delete(&book).Error
	if err != nil {
		return errors.New("BooksStorage error: " + err.Error())
	}
	return nil
}

func (storage *BooksStorage) GetAllBooks(limit, offset uint64) ([]model.Book, error) {
	var booksList []model.Book
	var err error

	if limit > 0 || offset > 0 {
		err = storage.db.Debug().Offset(offset).Limit(limit).Find(&booksList).Error
	} else {
		err = storage.db.Find(&booksList).Error
	}

	if err != nil {
		return []model.Book{}, errors.New("BooksStorage error: " + err.Error())
	}
	return booksList, nil
}

func (storage *BooksStorage) GetTopBooks(field string, limit uint64, asc bool) ([]model.Book, error) {
	if !model.IsBookField(field) {
		return []model.Book{}, errors.New("field " + field + " is not field of Book")
	}

	query := field
	if asc {
		query = query + " asc"
	} else {
		query = query + " desc"
	}

	var booksList []model.Book
	err := storage.db.Order(query).Limit(limit).Find(&booksList).Error
	if err != nil {
		return []model.Book{}, errors.New("BooksStorage error: " + err.Error())
	}
	return booksList, nil
}

func (storage *BooksStorage) SearchBooks(query string, distance, limit uint64) (booksList []model.Book, err error) {
	err = storage.db.Where("levenshtein(name, ?) <= ? OR levenshtein(author, ?) <= ?", query, distance, query, distance).
		Limit(limit).
		Find(&booksList).Error

	if err != nil {
		err = errors.New("BooksStorage error: " + err.Error())
	}
	return
}

func (storage *BooksStorage) GetFullBook(id uint64) (model.Book, error) {
	var book model.Book
	err := storage.db.Preload("Reviews").First(&book, "id = ?", id).Error

	if err != nil {
		return model.Book{}, errors.New("BooksStorage error: " + err.Error())
	}
	return book, nil
}
