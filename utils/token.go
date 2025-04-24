package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hackshel/tracker-server/pkg/setting"
)

var jwtKey = []byte(setting.JwtSecret)

type TokenClaims struct {
	Username string `json:"username"`
	UserRole string `json:"user_role"`
	UserID   string `json:"user_id"`
	Passkey  string `json:"passkey"`
	jwt.RegisteredClaims
}

func GenerateToken(username string, userRole string, userID string, Passkey string) (string, error) {
	fmt.Printf("userRole: %v\n", userRole)
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &TokenClaims{
		Username: username,
		UserRole: userRole,
		UserID:   userID,
		Passkey:  Passkey,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    setting.AppName,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func ParseToken(tokenStr string) (*TokenClaims, error) {
	claims := &TokenClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
