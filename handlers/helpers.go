package handlers

import (
	"bookstore/logger"
	"bookstore/model"
	"time"

	"net/http"
	"runtime"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

const messageKeyWord = "message"

func createErrorResponse(ctx *gin.Context, code int, errorMsg string) {
	logger.LogErrorResponce(code, errorMsg).Info("Creating error response")

	ctx.JSON(code, gin.H{messageKeyWord: errorMsg})
}

func createInternalError(ctx *gin.Context) {
	pc, _, _, _ := runtime.Caller(1)
	fullName := runtime.FuncForPC(pc).Name()
	parts := strings.Split(fullName, ".")
	createErrorResponse(ctx, http.StatusInternalServerError, "Internal error at "+parts[len(parts)-2])
}

func parseId(ctx *gin.Context) (id uint64, err error) {
	idString := ctx.Param("id")
	id, err = strconv.ParseUint(idString, 10, 64)
	if err != nil {
		createErrorResponse(ctx, http.StatusBadRequest, "Cannot parse id "+idString)
	}
	return
}

func parseFullFlag(ctx *gin.Context) (isFull bool, err error) {
	isFullString := ctx.DefaultQuery("full", "false")
	isFull, err = strconv.ParseBool(isFullString)
	if err != nil {
		createErrorResponse(ctx, http.StatusBadRequest, "Cannot parse full "+isFullString)
	}
	return
}

func hashPassword(ctx *gin.Context, password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		createErrorResponse(ctx, http.StatusInternalServerError, "Cannot hash password")
	}
	return string(bytes), err
}

func validatePassword(ctx *gin.Context, hashed, given string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(given))
	if err != nil {
		createErrorResponse(ctx, http.StatusUnauthorized, "Invalid username or password")
	}
	return err == nil, err
}

func generateAllTokens(user *model.User, secretKey []byte) (accessToken string, refreshToken string, err error) {
	// Create the access token
	accessTokenClaims := &model.JwtTokenClaims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Second * 10).Unix(),
		},
	}
	accessTokenToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessToken, err = accessTokenToken.SignedString(secretKey)
	if err != nil {
		return "", "", err
	}

	// Create the refresh token
	refreshTokenClaims := &model.JwtTokenClaims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 7).Unix(),
		},
	}

	refreshTokenToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	refreshToken, err = refreshTokenToken.SignedString(secretKey)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func getToken(ctx *gin.Context) (token *model.JwtTokenClaims, exists bool) {
	tokenAny, exists := ctx.Get("claims")
	if !exists {
		createErrorResponse(ctx, http.StatusUnauthorized, "No token found")
		return
	}

	token, ok := tokenAny.(*model.JwtTokenClaims)
	if !ok {
		createErrorResponse(ctx, http.StatusInternalServerError, "Cannot cast to model.JwtTokenClaims")
		return
	}

	return
}

func validatePermissions(ctx *gin.Context, requiredID uint64, requiredUsername, requiredRole string) bool {
	token, exist := getToken(ctx)
	if !exist {
		return false
	}

	if requiredID != 0 && token.UserID != requiredID {
		createErrorResponse(ctx, http.StatusForbidden, "Access denied for user id "+strconv.FormatUint(token.UserID, 10))
		return false
	}

	if len(requiredRole) != 0 && token.Role != requiredRole {
		createErrorResponse(ctx, http.StatusForbidden, "Access denied for user role "+token.Role)
		return false
	}

	if len(requiredUsername) != 0 && token.Username != requiredUsername {
		createErrorResponse(ctx, http.StatusForbidden, "Access denied for username "+token.Username)
		return false
	}

	return true
}

func refreshToken(refreshToken *model.JwtTokenClaims, secretKey []byte) (string, error) {
	accessTokenClaims := &model.JwtTokenClaims{
		UserID:   refreshToken.UserID,
		Username: refreshToken.Username,
		Role:     refreshToken.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Second * 10).Unix(),
		},
	}
	accessTokenToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessToken, err := accessTokenToken.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func parseLimitAndOffset(ctx *gin.Context) (limit, offset uint64, ok bool) {
	limitString := ctx.DefaultQuery("limit", "0")
	limit, err := strconv.ParseUint(limitString, 10, 64)
	if err != nil {
		createErrorResponse(ctx, http.StatusBadRequest, "Cannot parse limit "+limitString)
		return
	}

	offsetString := ctx.DefaultQuery("offset", "0")
	offset, err = strconv.ParseUint(offsetString, 10, 64)
	if err != nil {
		createErrorResponse(ctx, http.StatusBadRequest, "Cannot parse offset "+offsetString)
		return
	}

	ok = true

	return
}
