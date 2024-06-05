package middleware

import (
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type Claims struct {
	AccessToken string `json:"access_token"`
	jwt.StandardClaims
	UserId int64 `json:"id"`
}

func CreatToken(AccessToken string, id int64, secret string, expireTime int64) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		AccessToken: AccessToken,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Unix() + expireTime,
			Issuer:    "tiktok-app-microservice",
			Subject:   "access_token",
			IssuedAt:  time.Now().Unix(),
		},
		UserId: id,
	})
	token, err := claims.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return token, nil
}

func ValidToken(accessToken string, secret string) (bool, error) {
	tokenClaims, err := jwt.ParseWithClaims(accessToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			expireTime := claims.ExpiresAt
			nowTime := time.Now().Unix()
			if nowTime > expireTime {
				return true, nil
			} else {
				return false, nil
			}
		}
	}
	return true, err
}

func GetUserIdFromToken(accessToken string, secret string) (int64, error) {
	tokenClaims, err := jwt.ParseWithClaims(accessToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims.UserId, nil
		}
	}
	return -1, err
}
