package handlers

import (
	"bookstore/database"
	"bookstore/logger"
	"bookstore/model"

	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	db *database.UsersStorage
	secretKey []byte
}

func newUserHandler(config *Config) *UserHandler {
	return &UserHandler{
		db: config.UsersStorage,
		secretKey: config.SecretKey,
	}
}

func (handler *UserHandler) GetUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := parseId(ctx)
		if err != nil {
			return
		}

		if !validatePermissions(ctx, id, "", "") && !validatePermissions(ctx, 0, "", "admin") {
			return
		}

		isFull, err := parseFullFlag(ctx)
		if err != nil {
			return
		}

		getUser := func() (func(id uint64) (model.User, error)) {
			if isFull {
				return handler.db.GetFullUser
			} else {
				return handler.db.GetUser
			}
		}

		user, err := getUser()(id)
		if err != nil {
			logger.Log.Error("Error from storage: " + err.Error())

			if database.IsNotFoundError(err.Error()) {
				createErrorResponse(ctx, http.StatusNotFound, "Cannot find user with id "+strconv.FormatUint(id, 10))
			} else {
				createInternalError(ctx)
			}
			return
		}
		ctx.JSON(http.StatusOK, user)
	}
}

func (handler *UserHandler) GetAllUsers() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !validatePermissions(ctx, 0, "", "admin") {
			return
		}

		limit, offset, ok := parseLimitAndOffset(ctx)
		if !ok {
			logger.Log.Info("Not ok")
			return
		}

		users, err := handler.db.GetAllUsers(limit, offset)
		if err != nil {
			logger.Log.Error("Error from storage: " + err.Error())
			createInternalError(ctx)
			return
		}
		ctx.JSON(http.StatusOK, users)
	}
}

func parseUser(ctx *gin.Context) (user model.User, err error) {
	err = nil

	var short model.UserShort

	err = ctx.ShouldBindJSON(&short)
	if err != nil {
		return
	}

	user.Username = short.Username
	user.Password = short.Password
	user.Email = short.Email
	user.Role = short.Role

	return
}

func (handler *UserHandler) AddUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !validatePermissions(ctx, 0, "", "admin") {
			return
		}

		user, err := parseUser(ctx)
		if err != nil {
			logger.Log.Error("Error of parsing user: " + err.Error())
			createErrorResponse(ctx, http.StatusBadRequest, "Cannot parse user input")
			return
		}

		err = handler.db.AddUser(&user)
		if err != nil {
			logger.Log.Error("Error from storage: " + err.Error())
			createInternalError(ctx)
			return
		}

		ctx.JSON(http.StatusOK, user)
	}
}

func (handler *UserHandler) DeleteUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !validatePermissions(ctx, 0, "", "admin") {
			return
		}

		id, err := parseId(ctx)
		if err != nil {
			return
		}
		err = handler.db.DeleteUser(id)
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

func (handler *UserHandler) GetUserByUsername() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		username := ctx.Query("username")

		if !validatePermissions(ctx, 0, username, "") && !validatePermissions(ctx, 0, "", "admin") {
			return
		}

		user, err := handler.db.GetUserByUsername(username)
		if err != nil {
			logger.Log.Error("Error from storage: " + err.Error())

			if database.IsNotFoundError(err.Error()) {
				createErrorResponse(ctx, http.StatusNotFound, "Cannot find user with username "+username)
			} else {
				createInternalError(ctx)
			}
			return
		}
		user.Password = ""
		ctx.JSON(http.StatusOK, user)
	}
}

func (handler *UserHandler) SignUp() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, err := parseUser(ctx)
		if err != nil {
			logger.Log.Error("Error of parsing user: " + err.Error())
			createErrorResponse(ctx, http.StatusBadRequest, "Cannot parse user input")
			return
		}

		user.Role = "user"

		user.Password, err = hashPassword(ctx, user.Password)
		if err != nil {
			logger.Log.Error("Hashing error: " + err.Error())
			return
		}

		err = handler.db.AddUser(&user)
		if err != nil {
			logger.Log.Error("Error from storage: " + err.Error())
			createInternalError(ctx)
			return
		}

		accessToken, refreshToken, err := generateAllTokens(&user, handler.secretKey)
		if err != nil {
			logger.Log.Error("JWT tokens generator error: " + err.Error())
			createErrorResponse(ctx, http.StatusInternalServerError, "Cannot generate JWT tokens")
		}

		ctx.JSON(http.StatusOK, gin.H{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		})
	}
}

func (handler *UserHandler) Login() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var loginData struct {
			Username string `json:"username" binding:"required"`
			Password string `json:"password" binding:"required"`
		}
		if err := ctx.ShouldBindJSON(&loginData); err != nil {
			createErrorResponse(ctx, http.StatusInternalServerError, err.Error())
			return
		}

		user, err := handler.db.GetUserByUsername(loginData.Username)
		if err != nil {
			createErrorResponse(ctx, http.StatusUnauthorized, "Invalid username or password")
			return
		}

		isValid, err := validatePassword(ctx, user.Password, loginData.Password)
		if err != nil || !isValid {
			return
		}

		user.AccessToken, user.RefreshToken, err = generateAllTokens(&user, handler.secretKey)
		if err != nil {
			logger.Log.Error("JWT tokens generator error: " + err.Error())
			createErrorResponse(ctx, http.StatusInternalServerError, "Cannot generate JWT tokens")
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"access_token":  user.AccessToken,
			"refresh_token": user.RefreshToken,
		})
	}
}

func (handler *UserHandler) RefreshToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, exist := getToken(ctx)
		if !exist {
			return
		}

		newAccessToken, err := refreshToken(token, handler.secretKey)
		if err != nil {
			logger.Log.Error("Error from token refresh: " + err.Error())
			createErrorResponse(ctx, http.StatusInternalServerError, "Internal sever error while token refresh")
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"access_token":  newAccessToken,
		})
	}
}
