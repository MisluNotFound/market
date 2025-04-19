package lib

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mislu/market-api/internal/utils/app"
)

// TODO use file
var secretKey = "TODO USE FILE"

type CustomClaims struct {
	UserID        string `json:"userID"`
	IsAccessToken bool   `json:"isAccessToken"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(userID string) (string, error) {
	return generateToken(userID, time.Duration(app.GetConfig().Server.AccessTokenExpire), true)
}

func GenerateRefreshToken(userID string) (string, error) {
	return generateToken(userID, time.Duration(app.GetConfig().Server.RefreshTokenExpire), false)
}

func generateToken(userID string, expire time.Duration, isAccessToken bool) (string, error) {
	claims := CustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

func VerifyToken(tokenString string, isAccessToken bool) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid && claims.IsAccessToken == isAccessToken {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func refreshToken(refreshToken string) (string, error) {
	claims, err := VerifyToken(refreshToken, false)
	if err != nil {
		return "", fmt.Errorf("invalid refresh token: %v", err)
	}

	newAccessToken, err := GenerateAccessToken(claims.UserID)
	if err != nil {
		return "", fmt.Errorf("failed to generate new access token: %v", err)
	}

	return newAccessToken, nil
}
