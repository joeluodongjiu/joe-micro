package middleware

import (
	"github.com/gin-gonic/gin"
	"joe-micro/lib/jwt"
	"joe-micro/lib/log"
	"net/http"
)

// JWTAuth 中间件，检查token
func JWTAuth(skipper ...SkipperFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if len(skipper) > 0 && skipper[0](c) {
			c.Next()
			return
		}
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
		claims, err :=  jwt.ParseToken(token)
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