package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"joe-micro/adminApi/model/casbin"
	"joe-micro/lib/log"
	"strconv"
)

// CasbinMiddleware casbin中间件
func CasbinMiddleware(skipper ...SkipperFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if len(skipper) > 0 && skipper[0](c) {
			c.Next()
			return
		}
		var sub, obj, act string
		obj = c.Request.URL.Path
		method := c.Request.Method
		switch method {
		case "GET":
			act = "read"
		case "POST":
			act = "write"
		default:
			act = "read"
		}
		sub = strconv.Itoa(c.GetInt("uid"))
		log.Infof("权限认证:%v  %v  %v", sub, obj, act)
		if b, err := casbin.CsbinCheckPermission(sub, obj, act); err != nil {
			c.JSON(500, gin.H{
				"code": -1,
				"msg":  fmt.Sprintf("验证服务出现错误:%v", err),
			})
			c.Abort()
		} else if !b {
			c.JSON(400, gin.H{
				"code": 4,
				"msg":  "没有权限",
			})
			c.Abort()
		}
		c.Next()
	}
}
