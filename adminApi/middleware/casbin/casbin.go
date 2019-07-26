package casbin

import (
	"fmt"
	"github.com/casbin/casbin"
	"github.com/casbin/gorm-adapter"
	"github.com/gin-gonic/gin"
	"joe-micro/lib/log"
	"strconv"
)

var enforcer *casbin.Enforcer

func init() {
	/************************************/
	/********** casbin  权限管理  ********/
	/************************************/
	a := gormadapter.NewAdapter("mysql", "root:gogocuri@tcp(192.168.0.162:3306)/rbac", true)
	enforcer = casbin.NewEnforcer("./rbac.conf", a)
	//从DB加载策略
	err:=enforcer.LoadPolicy()
	if err!=nil{
		log.Fatal(err)
		return
	}
}

// CasbinMiddleware casbin中间件
func CasbinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var  sub,obj,act  string
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
		if b, err := enforcer.EnforceSafe(sub,obj,act); err != nil {
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
