package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/xweis/go-gin/pkg/setting"
	"github.com/xweis/go-gin/pkg/util/aes"
	"time"
)

var jwtKey = []byte(setting.AppSetting.JwtSecret)

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type TokenInfo struct {
	Token          string    `json:"token"`
	ExpirationTime time.Time `json:"expirationTime"`
}

func GetJwtToken(username string) (TokenInfo, error) {
	tokenInfo := TokenInfo{}

	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return tokenInfo, errors.New("failed to generate token")
	}

	//jwt aes 加密
	tokenInfo.Token, err = aes.EncryptByAes([]byte(tokenString))
	if err != nil {
		return TokenInfo{}, err
	}
	tokenInfo.ExpirationTime = expirationTime
	return tokenInfo, nil
}

func RefreshToken(tokenString string) (TokenInfo, error) {
	tokenInfo := TokenInfo{}
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return tokenInfo, err
	}

	if !token.Valid {
		return tokenInfo, errors.New("invalid token")
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	claims.ExpiresAt = jwt.NewNumericDate(expirationTime)
	token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(jwtKey)
	if err != nil {
		return tokenInfo, err
	}

	//jwt aes 加密
	tokenInfo.Token, err = aes.EncryptByAes([]byte(tokenString))
	if err != nil {
		return TokenInfo{}, err
	}
	tokenInfo.ExpirationTime = expirationTime
	return tokenInfo, nil
}

func ParseToken(token string) (*Claims, error) {
	tokenString, err := aes.DecryptByAes(token)
	if err != nil {
		return &Claims{}, err
	}
	tokenClaims, err := jwt.ParseWithClaims(string(tokenString), &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}
