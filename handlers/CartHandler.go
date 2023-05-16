package handlers

import (
	"bookstore/database"
	"bookstore/logger"
	"bookstore/model"

	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CartHandler struct {
	db        *database.CartsStorage
	secretKey []byte
}

func newCartHandler(config *Config) *CartHandler {
	return &CartHandler{
		db:        config.CartsStorage,
		secretKey: config.SecretKey,
	}
}

func (handler *CartHandler) GetCart() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := parseId(ctx)
		if err != nil {
			return
		}

		cart, err := handler.db.GetCart(id)
		if err != nil {
			logger.Log.Error("Error from storage: " + err.Error())

			if database.IsNotFoundError(err.Error()) {
				createErrorResponse(ctx, http.StatusNotFound, "Cannot find cart with id "+strconv.FormatUint(id, 10))
			} else {
				createInternalError(ctx)
			}
			return
		}

		if !validatePermissions(ctx, cart.UserID, "", "") && !validatePermissions(ctx, 0, "", "admin") {
			return
		}

		ctx.JSON(http.StatusOK, cart)
	}
}

func (handler *CartHandler) GetAllCarts() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !validatePermissions(ctx, 0, "", "admin") {
			return
		}

		limit, offset, ok := parseLimitAndOffset(ctx)
		if !ok {
			return
		}

		carts, err := handler.db.GetAllCarts(limit, offset)
		if err != nil {
			logger.Log.Error("Error from storage: " + err.Error())
			createInternalError(ctx)
			return
		}
		ctx.JSON(http.StatusOK, carts)
	}
}

func parseCart(ctx *gin.Context) (cart model.Cart, err error) {
	err = nil

	var short model.CartShort

	err = ctx.ShouldBindJSON(&short)
	if err != nil {
		return
	}

	cart.UserID = short.UserID
	cart.BookID = short.BookID

	return
}

func (handler *CartHandler) AddCart() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cart, err := parseCart(ctx)
		if err != nil {
			logger.Log.Error("Error of parsing cart: " + err.Error())
			createErrorResponse(ctx, http.StatusBadRequest, "Cannot parse cart input")
			return
		}

		if !validatePermissions(ctx, cart.UserID, "", "") && !validatePermissions(ctx, 0, "", "admin") {
			return
		}

		err = handler.db.AddCart(&cart)
		if err != nil {
			logger.Log.Error("Error from storage: " + err.Error())
			createInternalError(ctx)
			return
		}

		ctx.JSON(http.StatusOK, cart)
	}
}

func (handler *CartHandler) DeleteCart() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := parseId(ctx)
		if err != nil {
			return
		}

		cartFromDb, err := handler.db.GetCart(id)
		if err != nil {
			logger.Log.Error("Error from storage: " + err.Error())

			if database.IsNotFoundError(err.Error()) {
				createErrorResponse(ctx, http.StatusNotFound, "Cannot find book with id "+strconv.FormatUint(id, 10))
			} else {
				createInternalError(ctx)
			}
			return
		}

		if !validatePermissions(ctx, cartFromDb.UserID, "", "") && !validatePermissions(ctx, 0, "", "admin") {
			return
		}

		err = handler.db.DeleteCart(id)
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

func (handler *CartHandler) GetCartsByUserId() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userIdString := ctx.Query("userid")
		userId, err := strconv.ParseUint(userIdString, 10, 64)
		if err != nil {
			createErrorResponse(ctx, http.StatusBadRequest, "Cannot parse user id "+userIdString)
			return
		}

		if !validatePermissions(ctx, userId, "", "") && !validatePermissions(ctx, 0, "", "admin") {
			return
		}

		limit, offset, ok := parseLimitAndOffset(ctx)
		if !ok {
			return
		}

		carts, err := handler.db.GetCartsByUserId(userId, limit, offset)
		if err != nil {
			logger.Log.Error("Error from storage: " + err.Error())
			createInternalError(ctx)
			return
		}
		ctx.JSON(http.StatusOK, carts)
	}
}

func (handler *CartHandler) GetCartsByBookId() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !validatePermissions(ctx, 0, "", "admin") {
			return
		}

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

		carts, err := handler.db.GetCartsByBookId(bookId, limit, offset)
		if err != nil {
			logger.Log.Error("Error from storage: " + err.Error())
			createInternalError(ctx)
			return
		}
		ctx.JSON(http.StatusOK, carts)
	}
}
