package routers

import (
	"github.com/gin-gonic/gin"
	"joe-micro/adminApi/handler"
	"joe-micro/adminApi/middleware/casbin"
	"joe-micro/lib/jwt"
	"joe-micro/lib/log"
	"joe-micro/lib/trace"
)



func Init()  *gin.Engine{
	gin.SetMode(gin.ReleaseMode) //是否生产模式启动
	router := gin.Default()
	router.Use(log.GinLogger())
	router.Use(trace.TracerWrapper)
	// 跨域
	router.Use(func(ctx *gin.Context) {
		ctx.Header("Access-Control-Allow-Origin", "*")                            //跨域
		ctx.Header("Access-Control-Allow-Headers", "Token,Content-Type")          //必须的请求头
		ctx.Header("Access-Control-Allow-Methods", "OPTIONS,PUT,POST,GET,DELETE") //接收的请求方法
	})

	// OPTIONS
	router.NoRoute(func(ctx *gin.Context) {
		if ctx.Request.Method == "OPTIONS" {
			ctx.JSON(200, nil)
		}
	})

	//登录接口
	router.POST("/admin/login", handler.Login)

	//jwt 用户鉴权
	router.Use(jwt.JWTAuth())

	//casbin  权限管理
	router.Use(casbin.CasbinMiddleware())

	//swagger
	/*	url := ginSwagger.URL("http://localhost:9081/swagger/doc.json") // The url pointing to API definition
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))*/

	admin := router.Group("/admin")
	admin.POST("/login/exit", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"code": 0,
			"msg":  "ok",
		})
		return
	})


	user := admin.Group("/user")
	{
		user.POST("createOne", handler.CreateOne)
	}


	permission := admin.Group("/permission")
	{
		permission.GET("/readResource", func(c *gin.Context) {
			log.Info("canRead")
			c.JSON(200, gin.H{
				"code": 0,
				"msg":  "ok",
			})
			return
		})
		permission.POST("/writeResource", func(c *gin.Context) {
			log.Info("canWrite")
			c.JSON(200, gin.H{
				"code": 0,
				"msg":  "ok",
			})
			return
		})

	}
    return router
}