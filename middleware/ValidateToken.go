package middleware

import (
	"bookstore/model"

	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func ValidateToken(secretKey []byte) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get token from Authorization header
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		// Parse token
		token, err := jwt.ParseWithClaims(tokenString, &model.JwtTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token signature"})
				return
			}
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		if !token.Valid {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token is not valid"})
			return
		}

		// Save token claims to context
		claims, ok := token.Claims.(*model.JwtTokenClaims)
		if ok && token.Valid {
			ctx.Set("claims", claims)
		} else {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		ctx.Next()
	}
}
