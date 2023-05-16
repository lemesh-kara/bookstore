package model

import "github.com/golang-jwt/jwt"

type JwtTokenClaims struct {
	UserID   uint64 `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}
