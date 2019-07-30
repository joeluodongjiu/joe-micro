package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	_ "joe-micro/adminApi/docs"
	"joe-micro/adminApi/handler"
	"joe-micro/adminApi/middleware"
	"joe-micro/lib/log"
	"joe-micro/lib/trace"
)

func Init() *gin.Engine {
	gin.SetMode(gin.ReleaseMode) //是否生产模式启动
	router := gin.Default()

	router.NoRoute(middleware.NoRouteHandler())
	// 崩溃恢复
	router.Use(middleware.RecoveryMiddleware())
	// gin日志
	router.Use(log.GinLogger())

	//swagger
	url := ginSwagger.URL("http://localhost:9081/swagger/doc.json") // The url pointing to API definition
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	// jaeger trace 追踪
	router.Use(trace.TracerWrapper)
	// 跨域
	router.Use(func(ctx *gin.Context) {
		ctx.Header("Access-Control-Allow-Origin", "*")                   //跨域
		ctx.Header("Access-Control-Allow-Headers", "token,Content-Type") //必须的请求头
		ctx.Header("Access-Control-Allow-Methods", "OPTIONS,POST,GET")   //接收的请求方法
	})

	RegisterRouter(router)
	return router
}

func RegisterRouter(api *gin.Engine) {
	apiPrefix := "/api/admin"

	// 登录验证 jwt token 验证 及信息提取
	var notCheckLoginUrlArr []string
	notCheckLoginUrlArr = append(notCheckLoginUrlArr, apiPrefix+"/user/login")
	//jwt 用户鉴权
	api.Use(middleware.JWTAuth(middleware.AllowPathPrefixSkipper(notCheckLoginUrlArr...)))

	//casbin  权限管理
	var notCheckPermissionUrlArr []string
	notCheckPermissionUrlArr = append(notCheckPermissionUrlArr, notCheckLoginUrlArr...)
	notCheckPermissionUrlArr = append(notCheckPermissionUrlArr, apiPrefix+"/user/logout")
	api.Use(middleware.CasbinMiddleware(middleware.AllowPathPrefixSkipper(notCheckPermissionUrlArr...)))

	admin := api.Group(apiPrefix)
	user := handler.User{} //用户模块
	admin.POST("/user/login", user.Login)
	admin.GET("/user/logout", user.Logout)
	admin.GET("/user/info", user.Info)
	admin.POST("/user/edit_pwd", user.EditPwd)

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
}
