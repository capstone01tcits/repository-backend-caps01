package utils

import (
	"errors"
	"strconv"
	"time"

	"go-auth/config"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Type   string `json:"type"` // "access" or "refresh"
	jwt.RegisteredClaims
}

func GenerateAccessToken(userID, email string) (string, int64, error) {
	hours, _ := strconv.Atoi(config.Cfg.JWTExpireHours)
	expireTime := time.Now().Add(time.Duration(hours) * time.Hour)

	claims := &Claims{
		UserID: userID,
		Email:  email,
		Type:   "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(config.Cfg.JWTSecret))
	return signed, expireTime.Unix(), err
}

func GenerateRefreshToken(userID, email string) (string, error) {
	hours, _ := strconv.Atoi(config.Cfg.JWTRefreshExpireHours)
	expireTime := time.Now().Add(time.Duration(hours) * time.Hour)

	claims := &Claims{
		UserID: userID,
		Email:  email,
		Type:   "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Cfg.JWTRefreshSecret))
}

func ValidateAccessToken(tokenStr string) (*Claims, error) {
	return validateToken(tokenStr, config.Cfg.JWTSecret, "access")
}

func ValidateRefreshToken(tokenStr string) (*Claims, error) {
	return validateToken(tokenStr, config.Cfg.JWTRefreshSecret, "refresh")
}

func validateToken(tokenStr, secret, expectedType string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	if claims.Type != expectedType {
		return nil, errors.New("wrong token type")
	}

	return claims, nil
}
