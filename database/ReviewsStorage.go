package database

import (
	"bookstore/model"
	"errors"

	"github.com/jinzhu/gorm"
)

type ReviewsStorage struct {
	db *gorm.DB
}

func NewReviewsStorage(db *gorm.DB) *ReviewsStorage {
	return &ReviewsStorage{
		db: db,
	}
}

func (storage *ReviewsStorage) AddReview(review *model.Review) error {
	prevReviews, err := storage.GetReviewsByBookId(review.BookID, 0, 0)
	if err != nil {
		return errors.New("ReviewsStorage error: " + err.Error())
	}

	newMark := float64(review.ReviewMark)
	amountOfMarks := len(prevReviews) + 1
	for _, prevReview := range prevReviews {
		newMark += prevReview.ReviewMark
	}
	newMark /= float64(amountOfMarks)

	var book model.Book
	err = storage.db.First(&book, "id = ?", review.BookID).Error
	if err != nil {
		return errors.New("ReviewsStorage error: " + err.Error())
	}

	book.ReviewMark = newMark
	err = storage.db.Save(&book).Error
	if err != nil {
		return errors.New("ReviewsStorage error: " + err.Error())
	}

	err = storage.db.Create(review).Error
	if err != nil {
		return errors.New("ReviewsStorage error: " + err.Error())
	}
	return nil
}

func (storage *ReviewsStorage) GetReview(id uint64) (model.Review, error) {
	var review model.Review
	err := storage.db.Preload("User").Preload("Book").First(&review, "id = ?", id).Error

	if err != nil {
		return model.Review{}, errors.New("ReviewsStorage error: " + err.Error())
	}

	resetUserSecrets(&review.User)

	return review, nil
}

func (storage *ReviewsStorage) DeleteReview(id uint64) error {
	review, err := storage.GetReview(id)
	if err != nil {
		return errors.New("ReviewsStorage error: " + err.Error())
	}
	err = storage.db.Delete(&review).Error
	if err != nil {
		return errors.New("ReviewsStorage error: " + err.Error())
	}

	prevReviews, err := storage.GetReviewsByBookId(review.BookID, 0, 0)
	if err != nil {
		return errors.New("ReviewsStorage error: " + err.Error())
	}

	newMark := float64(0)
	amountOfMarks := len(prevReviews)

	if amountOfMarks == 0 {
		return nil
	}

	for _, prevReview := range prevReviews {
		newMark += prevReview.ReviewMark
	}
	newMark /= float64(amountOfMarks)

	var book model.Book
	err = storage.db.First(&book, "id = ?", review.BookID).Error
	if err != nil {
		return errors.New("ReviewsStorage error: " + err.Error())
	}

	book.ReviewMark = newMark
	err = storage.db.Save(&book).Error
	if err != nil {
		return errors.New("ReviewsStorage error: " + err.Error())
	}
	return nil
}

func (storage *ReviewsStorage) GetAllReviews(limit, offset uint64) ([]model.Review, error) {
	var reviewsList []model.Review
	var err error

	if limit > 0 || offset > 0 {
		err = storage.db.Offset(offset).Limit(limit).Find(&reviewsList).Error
	} else {
		err = storage.db.Find(&reviewsList).Error
	}

	if err != nil {
		return []model.Review{}, errors.New("ReviewsStorage error: " + err.Error())
	}
	return reviewsList, nil
}

func (storage *ReviewsStorage) GetReviewsByUserId(userID uint64, limit, offset uint64) ([]model.Review, error) {
	var reviews []model.Review
	var err error

	if limit > 0 || offset > 0 {
		err = storage.db.Offset(offset).Limit(limit).Preload("User").Preload("Book").Where("user_id = ?", userID).Find(&reviews).Error
	} else {
		err = storage.db.Preload("User").Preload("Book").Where("user_id = ?", userID).Find(&reviews).Error
	}

	if err != nil {
		return []model.Review{}, errors.New("ReviewsStorage error: " + err.Error())
	}

	for index, _ := range reviews {
		resetUserSecrets(&reviews[index].User)
	}

	return reviews, nil
}

func (storage *ReviewsStorage) GetReviewsByBookId(bookID uint64, limit, offset uint64) ([]model.Review, error) {
	var reviews []model.Review
	var err error

	if limit > 0 || offset > 0 {
		err = storage.db.Offset(offset).Limit(limit).Preload("User").Preload("Book").Where("book_id = ?", bookID).Find(&reviews).Error
	} else {
		err = storage.db.Preload("User").Preload("Book").Where("book_id = ?", bookID).Find(&reviews).Error
	}

	if err != nil {
		return []model.Review{}, errors.New("ReviewsStorage error: " + err.Error())
	}

	for index, _ := range reviews {
		resetUserSecrets(&reviews[index].User)
	}

	return reviews, nil
}
