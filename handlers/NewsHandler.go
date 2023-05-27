package handlers

import (
	"bookstore/database"
	"bookstore/logger"

	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type NewsHandler struct {
	db *database.NewsStorage
}

func newNewsHandler(config *Config) *NewsHandler {
	return &NewsHandler{db: config.NewsStorage}
}

func (handler *NewsHandler) GetNews() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := parseId(ctx)
		if err != nil {
			return
		}

		cart, err := handler.db.GetNews(id)
		if err != nil {
			logger.Log.Error("Error from storage: " + err.Error())

			if database.IsNotFoundError(err.Error()) {
				createErrorResponse(ctx, http.StatusNotFound, "Cannot find news with id "+strconv.FormatUint(id, 10))
			} else {
				createInternalError(ctx)
			}
			return
		}

		ctx.JSON(http.StatusOK, cart)
	}
}

func (handler *NewsHandler) GetAllNews() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		news, err := handler.db.GetAllNews()
		if err != nil {
			logger.Log.Error("Error from storage: " + err.Error())
			createInternalError(ctx)
			return
		}
		ctx.JSON(http.StatusOK, news)
	}
}
