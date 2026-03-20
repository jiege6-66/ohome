package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
)

type JwtClaims struct {
	ID   uint
	Name string
	jwt.RegisteredClaims
}

var stSignKey = []byte(viper.GetString("jwt.signKey"))

// GenerateAccessToken 获取accessToken
func GenerateAccessToken(id uint, name string) (string, error) {
	token, err := generateToken(id, name, viper.GetDuration("jwt.accessTokenExpires")*time.Minute)
	return token, err
}

// GenerateRefreshToken 获取refreshToken
func GenerateRefreshToken(id uint, name string) (string, error) {
	token, err := generateToken(id, name, viper.GetDuration("jwt.refreshTokenExpires")*time.Minute)
	return token, err
}

func generateToken(id uint, name string, expiresTime time.Duration) (string, error) {
	iJwtCustClaims := JwtClaims{
		ID:   id,
		Name: name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   "Token",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, iJwtCustClaims)
	return token.SignedString(stSignKey)
}

func ParseToken(tokenStr string) (JwtClaims, error) {
	iJwtCustClaims := JwtClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, &iJwtCustClaims, func(token *jwt.Token) (interface{}, error) {
		return stSignKey, nil
	})

	if err == nil && !token.Valid {
		err = errors.New("Invalid Token")
	}

	return iJwtCustClaims, err
}
