package routers

import (
	"github.com/gin-gonic/gin"
	//"github.com/swaggo/files"
	//"github.com/swaggo/gin-swagger"
	_ "joe-micro/adminApi/docs"
	"joe-micro/adminApi/handler"
	"joe-micro/adminApi/middleware"
	"joe-micro/lib/log"
	"joe-micro/lib/trace"
	"joe-micro/lib/validator"
)

func Init() *gin.Engine {
	gin.SetMode(gin.ReleaseMode) //是否生产模式启动
	validator.Init()             //自定义验证器 替换 gin 里面的验证器
	router := gin.Default()

	router.NoRoute(middleware.NoRouteHandler())
	// 崩溃恢复
	router.Use(middleware.RecoveryMiddleware())
	// gin日志
	router.Use(log.GinLogger())

	//swagger
	/*	url := ginSwagger.URL("http://localhost:9081/swagger/doc.json") // The url pointing to API definition
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))*/

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
	admin := api.Group(apiPrefix)

	// 登录验证 jwt token 验证 及信息提取
	var notCheckLoginUrlArr []string
	notCheckLoginUrlArr = append(notCheckLoginUrlArr, apiPrefix+"/user/login")
	//jwt 用户鉴权
	admin.Use(middleware.JWTAuth(middleware.AllowPathPrefixSkipper(notCheckLoginUrlArr...)))

	//casbin  权限管理
	var notCheckPermissionUrlArr []string
	notCheckPermissionUrlArr = append(notCheckPermissionUrlArr, notCheckLoginUrlArr...)
	notCheckPermissionUrlArr = append(notCheckPermissionUrlArr, apiPrefix+"/user/logout")
	notCheckPermissionUrlArr = append(notCheckPermissionUrlArr, apiPrefix+"/user/edit_pwd")
	notCheckPermissionUrlArr = append(notCheckPermissionUrlArr, apiPrefix+"/user/info")
	admin.Use(middleware.CasbinMiddleware(middleware.AllowPathPrefixSkipper(notCheckPermissionUrlArr...)))
	//用户操作
	userC := handler.UserController{}
	admin.POST("/user/login", userC.Login)
	admin.GET("/user/logout", userC.Logout)
	admin.GET("/user/info", userC.Info)
	admin.POST("/user/edit_pwd", userC.EditPwd)
	//用户管理
	user_manaC := handler.UserManagementController{}
	admin.GET("/user_mana/list", user_manaC.List)
	admin.GET("/user_mana/detail", user_manaC.Detail)
	admin.POST("/user_mana/delete", user_manaC.Delete)
	admin.POST("/user_mana/update", user_manaC.Update)
	admin.POST("/user_mana/create", user_manaC.Create)
	admin.GET("/user_mana/users_roleid_list", user_manaC.UsersRoleIDList)
	admin.POST("/user_mana/set_role", user_manaC.SetRole)
	//角色管理
	roleC := handler.RoleController{}
	admin.GET("/role/list", roleC.List)
	admin.GET("/role/detail", roleC.Detail)
	admin.POST("/role/update", roleC.Update)
	admin.POST("/role/delete", roleC.Delete)
	admin.POST("/role/create", roleC.Create)
	admin.GET("/role/role_menuid_list", roleC.RoleMenuIDList)
	admin.POST("/role/set_role_with_menus", roleC.SetRoleWithMenus)
    //菜单管理
	menuC := handler.MenuController{}
	admin.GET("/menu/list", menuC.List)
	admin.GET("/menu/detail", menuC.Detail)
	admin.POST("/menu/delete", menuC.Delete)
	admin.POST("/menu/update", menuC.Update)
	admin.POST("/menu/create", menuC.Create)
	admin.GET("/menu/allmenu", menuC.AllMenu)
	admin.GET("/menu/menubuttonlist", menuC.MenuButtonList)
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
