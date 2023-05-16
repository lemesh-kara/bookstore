package handlers

import (
	"bookstore/database"
	"bookstore/logger"
	"bookstore/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BookHandler struct {
	db        *database.BooksStorage
	secretKey []byte
}

func newBookHandler(config *Config) *BookHandler {
	return &BookHandler{
		db:        config.BooksStorage,
		secretKey: config.SecretKey,
	}
}

func (bh *BookHandler) GetBook() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idString := ctx.Param("id")
		id, err := strconv.ParseUint(idString, 10, 64)
		if err != nil {
			createErrorResponse(ctx, http.StatusBadRequest, "Cannot parse id "+idString)
			return
		}

		isFull, err := parseFullFlag(ctx)
		if err != nil {
			return
		}

		getBook := func() func(id uint64) (model.Book, error) {
			if isFull {
				return bh.db.GetFullBook
			} else {
				return bh.db.GetBook
			}
		}

		book, err := getBook()(id)
		if err != nil {
			logger.Log.Error("Error from storage: " + err.Error())

			if database.IsNotFoundError(err.Error()) {
				createErrorResponse(ctx, http.StatusNotFound, "Cannot find book with id "+idString)
			} else {
				createInternalError(ctx)
			}
			return
		}
		ctx.JSON(http.StatusOK, book)
	}
}

func (bh *BookHandler) GetAllBooks() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		limit, offset, ok := parseLimitAndOffset(ctx)
		if !ok {
			return
		}

		books, err := bh.db.GetAllBooks(limit, offset)
		if err != nil {
			logger.Log.Error("Error from storage: " + err.Error())
			createInternalError(ctx)
			return
		}
		ctx.JSON(http.StatusOK, books)
	}
}

type sortingParams struct {
	field string
	limit uint64
	isAsc bool
}

func parseSortingParams(ctx *gin.Context) (params sortingParams, status bool) {
	status = true
	var err error

	params.field = ctx.DefaultQuery("field", "review_mark")

	limitString := ctx.DefaultQuery("limit", "5")
	params.limit, err = strconv.ParseUint(limitString, 10, 64)
	if err != nil {
		createErrorResponse(ctx, http.StatusBadRequest, "Cannot parse limit "+limitString)
		status = false
	}

	ascString := ctx.DefaultQuery("asc", "true")
	params.isAsc, err = strconv.ParseBool(ascString)
	if err != nil {
		createErrorResponse(ctx, http.StatusBadRequest, "Cannot parse asc "+ascString)
		status = false
	}

	return
}

func (bh *BookHandler) GetTop() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		params, status := parseSortingParams(ctx)
		if !status {
			return
		}
		books, err := bh.db.GetTopBooks(params.field, params.limit, params.isAsc)
		if err != nil {
			logger.Log.Error("Error from storage: " + err.Error())
			createInternalError(ctx)
			return
		}
		ctx.JSON(http.StatusOK, books)
	}
}

type searchingParams struct {
	query    string
	distance uint64
	limit    uint64
}

func parseSearchingParams(ctx *gin.Context) (params searchingParams, status bool) {
	status = true
	var err error

	params.query = ctx.Query("query")

	distanceString := ctx.DefaultQuery("distance", "5")
	params.distance, err = strconv.ParseUint(distanceString, 10, 64)
	if err != nil {
		createErrorResponse(ctx, http.StatusBadRequest, "Cannot parse distance "+distanceString)
		status = false
	}

	limitString := ctx.DefaultQuery("limit", "5")
	params.limit, err = strconv.ParseUint(limitString, 10, 64)
	if err != nil {
		createErrorResponse(ctx, http.StatusBadRequest, "Cannot parse limit "+limitString)
		status = false
	}

	return
}

func (bh *BookHandler) SearchBooks() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		params, status := parseSearchingParams(ctx)
		if !status {
			return
		}

		books, err := bh.db.SearchBooks(params.query, params.distance, params.limit)
		if err != nil {
			logger.Log.Error("Error from storage: " + err.Error())
			createInternalError(ctx)
			return
		}
		ctx.JSON(http.StatusOK, books)
	}
}

func parseBook(ctx *gin.Context) (book model.Book, err error) {
	err = nil

	type ShortBook struct {
		Name        string  `json:"name"`
		Author      string  `json:"author"`
		Year        string  `json:"year"`
		Description string  `json:"description"`
		PathToPdf   string  `json:"path_to_pdf"`
		ReviewMark  float64 `json:"review_mark"`
	}

	var shortBook ShortBook

	err = ctx.ShouldBindJSON(&shortBook)
	if err != nil {
		return
	}

	book.Name = shortBook.Name
	book.Author = shortBook.Author
	book.Year = shortBook.Year
	book.Description = shortBook.Description
	book.PathToPdf = shortBook.PathToPdf
	book.ReviewMark = shortBook.ReviewMark

	return
}

func (bh *BookHandler) AddBook() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !validatePermissions(ctx, 0, "", "admin") {
			return
		}

		book, err := parseBook(ctx)
		if err != nil {
			logger.Log.Error("Internal error for parsing book: " + err.Error())
			createErrorResponse(ctx, http.StatusBadRequest, "Cannot parse book input")
			return
		}

		err = bh.db.AddBook(&book)
		if err != nil {
			logger.Log.Error("Error from storage: " + err.Error())
			createInternalError(ctx)
			return
		}

		ctx.JSON(http.StatusOK, book)
	}
}

func (bh *BookHandler) DeleteBook() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !validatePermissions(ctx, 0, "", "admin") {
			return
		}

		idString := ctx.Param("id")
		id, err := strconv.ParseUint(idString, 10, 64)
		if err != nil {
			createErrorResponse(ctx, http.StatusBadRequest, "Cannot parse id "+idString)
			return
		}
		err = bh.db.DeleteBook(id)
		if err != nil {
			logger.Log.Error("Error from storage: " + err.Error())

			if database.IsNotFoundError(err.Error()) {
				createErrorResponse(ctx, http.StatusNotFound, "Cannot find book with id "+idString)
			} else {
				createInternalError(ctx)
			}
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "removed"})
	}
}
