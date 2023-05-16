package handlers

import (
	"bookstore/database"
	"bookstore/logger"
	"bookstore/model"

	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ReviewHandler struct {
	db        *database.ReviewsStorage
	secretKey []byte
}

func newReviewHandler(config *Config) *ReviewHandler {
	return &ReviewHandler{
		db:        config.ReviewsStorage,
		secretKey: config.SecretKey,
	}
}

func (handler *ReviewHandler) GetReview() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := parseId(ctx)
		if err != nil {
			return
		}

		review, err := handler.db.GetReview(id)
		if err != nil {
			logger.Log.Error("Error from storage: " + err.Error())

			if database.IsNotFoundError(err.Error()) {
				createErrorResponse(ctx, http.StatusNotFound, "Cannot find review with id "+strconv.FormatUint(id, 10))
			} else {
				createInternalError(ctx)
			}
			return
		}
		ctx.JSON(http.StatusOK, review)
	}
}

func (handler *ReviewHandler) GetAllReviews() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		limit, offset, ok := parseLimitAndOffset(ctx)
		if !ok {
			return
		}

		reviews, err := handler.db.GetAllReviews(limit, offset)
		if err != nil {
			logger.Log.Error("Error from storage: " + err.Error())
			createInternalError(ctx)
			return
		}
		ctx.JSON(http.StatusOK, reviews)
	}
}

func parseReview(ctx *gin.Context) (review model.Review, err error) {
	err = nil

	var short model.ReviewShort

	err = ctx.ShouldBindJSON(&short)
	if err != nil {
		return
	}

	review.UserID = short.UserID
	review.BookID = short.BookID
	review.ReviewMark = short.ReviewMark
	review.Text = short.Text

	return
}

func (handler *ReviewHandler) AddReview() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		review, err := parseReview(ctx)
		if err != nil {
			logger.Log.Error("Error of parsing review: " + err.Error())
			createErrorResponse(ctx, http.StatusBadRequest, "Cannot parse review input")
			return
		}

		if !validatePermissions(ctx, review.UserID, "", "") && !validatePermissions(ctx, 0, "", "admin") {
			return
		}

		err = handler.db.AddReview(&review)
		if err != nil {
			logger.Log.Error("Error from storage: " + err.Error())
			createInternalError(ctx)
			return
		}

		ctx.JSON(http.StatusOK, review)
	}
}

func (handler *ReviewHandler) DeleteReview() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := parseId(ctx)
		if err != nil {
			return
		}
		if !validatePermissions(ctx, id, "", "") && !validatePermissions(ctx, 0, "", "admin") {
			return
		}

		err = handler.db.DeleteReview(id)
		if err != nil {
			logger.Log.Error("Error from storage: " + err.Error())

			if database.IsNotFoundError(err.Error()) {
				createErrorResponse(ctx, http.StatusNotFound, "Cannot find book with id "+strconv.FormatUint(id, 10))
			} else {
				createInternalError(ctx)
			}
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "removed"})
	}
}

func (handler *ReviewHandler) GetReviewsByUserId() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userIdString := ctx.Query("userid")
		userId, err := strconv.ParseUint(userIdString, 10, 64)
		if err != nil {
			createErrorResponse(ctx, http.StatusBadRequest, "Cannot parse user id "+userIdString)
			return
		}

		limit, offset, ok := parseLimitAndOffset(ctx)
		if !ok {
			return
		}

		reviews, err := handler.db.GetReviewsByUserId(userId, limit, offset)
		if err != nil {
			logger.Log.Error("Error from storage: " + err.Error())
			createInternalError(ctx)
			return
		}
		ctx.JSON(http.StatusOK, reviews)
	}
}

func (handler *ReviewHandler) GetReviewsByBookId() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		bookIdString := ctx.Query("bookid")
		bookId, err := strconv.ParseUint(bookIdString, 10, 64)
		if err != nil {
			createErrorResponse(ctx, http.StatusBadRequest, "Cannot parse book id "+bookIdString)
			return
		}

		limit, offset, ok := parseLimitAndOffset(ctx)
		if !ok {
			return
		}

		reviews, err := handler.db.GetReviewsByBookId(bookId, limit, offset)
		if err != nil {
			logger.Log.Error("Error from storage: " + err.Error())
			createInternalError(ctx)
			return
		}
		ctx.JSON(http.StatusOK, reviews)
	}
}
