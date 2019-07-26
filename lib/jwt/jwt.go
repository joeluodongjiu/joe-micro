package jwt

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"joe-micro/lib/log"
	"net/http"
	"time"
)

// 一些常量
var (
	TokenExpired       = errors.New("Token 已过期")
	TokenNotValidYet   = errors.New("Token 未激活")
	TokenMalformed     = errors.New("这不是 Token")
	TokenInvalid       = errors.New("无法解析的 Token")
	SignKey             = []byte("newtrekWang")
)

// JWTAuth 中间件，检查token
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("token")
		if token == "" {
			c.JSON(http.StatusOK, gin.H{
				"code":  4,
				"msg":  "请求未携带token，无权限访问",
			})
			c.Abort()
			return
		}
		log.Info("get token: ", token)
		// parseToken 解析token包含的信息
		claims, err :=  parseToken(token)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":  4,
				"msg":    err.Error(),
			})
			c.Abort()
			return
		}
		// 继续交由下一个路由处理,并将解析出的信息传递下去
		c.Set("uid", claims.UID)
	}
}



// 载荷，可以加一些自己需要的信息
type customClaims struct {
	UID    int `json:"uid"`
	jwt.StandardClaims
}






// CreateToken 生成一个token
func  CreateToken(uid int) (string, error) {
	claims := &customClaims{}
	claims.UID = uid
	claims.ExpiresAt= time.Now().Add(10*time.Hour).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(SignKey)
}

// 解析Tokne
func parseToken(tokenString string) (*customClaims, error) {
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
func  RefreshToken(tokenString string) (string, error) {
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
