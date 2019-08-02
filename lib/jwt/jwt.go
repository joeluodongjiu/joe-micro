package jwt

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"joe-micro/lib/config"
	"time"
)

// 一些常量
var (
	TokenExpired     = errors.New("Token 已过期")
	TokenNotValidYet = errors.New("Token 未激活")
	TokenMalformed   = errors.New("这不是 Token")
	TokenInvalid     = errors.New("无法解析的 Token")
	SignKey          = []byte(config.C.Jwt.SignKey)
)

// 载荷，可以加一些自己需要的信息
type customClaims struct {
	UID string `json:"uid"`
	jwt.StandardClaims
}

// CreateToken 生成一个token
func CreateToken(uid string) (string, error) {
	claims := &customClaims{}
	claims.UID = uid
	claims.ExpiresAt = time.Now().Add(config.C.Jwt.TimeOut * time.Hour).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(SignKey)
}

// 解析Tokne
func ParseToken(tokenString string) (*customClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &customClaims{}, func(token *jwt.Token) (interface{}, error) {
		return SignKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if claims, ok := token.Claims.(*customClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, TokenInvalid
}

// 更新token
func RefreshToken(tokenString string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt.ParseWithClaims(tokenString, &customClaims{}, func(token *jwt.Token) (interface{}, error) {
		return SignKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*customClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
		return CreateToken(claims.UID)
	}
	return "", TokenInvalid
}
