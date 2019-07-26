package main

import (
	"github.com/casbin/casbin"
	"github.com/casbin/gorm-adapter"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/consul"
	microLog "github.com/micro/go-micro/util/log"
	"github.com/micro/go-micro/web"
	"github.com/opentracing/opentracing-go"
	"joe-micro/api/handler"
	"joe-micro/lib/config"
	"joe-micro/lib/log"
	"joe-micro/lib/queue"
	"joe-micro/lib/trace"
	"time"
)

// @title  微服务的管理端api文档demo
// @version 1.0
// @host  localhost:9081
// @BasePath /
func main() {
	//统一日志到服务的日志
	microLog.SetLogger(log.NewMicroLogger())

	/************************************/
	/********** 服务发现  cousul   ********/
	/************************************/
	reg := consul.NewRegistry(func(op *registry.Options) {
		op.Addrs = []string{
			config.C.Consul,
		}

	})
	// create new api service
	service := web.NewService(
		web.Name(config.C.Service.Name),
		web.Registry(reg),
		web.RegisterTTL(time.Second*15),      //重新注册时间
		web.RegisterInterval(time.Second*10), //注册过期时间
		web.Version(config.C.Service.Version),
		web.Address(config.C.Service.Port),
	)

	// initialise service
	if err := service.Init(); err != nil {
		log.Error(err.Error())
	}

	/************************************/
	/********** 链路追踪  trace   ********/
	/************************************/
	trace.SetSamplingFrequency(50)
	t, io, err := trace.NewTracer(config.C.Service.Name, config.C.Jaeger)
	if err != nil {
		log.Fatal(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	/************************************/
	/********** 消息队列  queue   ********/
	/************************************/
	queue.Init(config.C.Nsq.Address, config.C.Nsq.Lookup, config.C.Nsq.MaxInFlight)

	/************************************/
	/********** casbin  权限管理  ********/
	/************************************/
	a := gormadapter.NewAdapter("mysql", "root:gogocuri@tcp(192.168.0.162:3306)/rbac", true)
	e := casbin.NewEnforcer("./rbac.conf", a)
	//从DB加载策略
	e.LoadPolicy()

	/************************************/
	/********** gin  路由框架     ********/
	/************************************/
	gin.SetMode(gin.ReleaseMode) //是否生产模式启动
	router := gin.Default()
	router.Use(log.GinLogger())
	router.Use(trace.TracerWrapper)

	router.Use(LanjieqiHandler(e))

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

	//swagger
	/*	url := ginSwagger.URL("http://localhost:9081/swagger/doc.json") // The url pointing to API definition
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))*/

	admin := router.Group("/admin")
	{
		admin.GET("")
	}

	// register html handler
	service.Handle("/", router)

	// register call handler

	service.HandleFunc("/api/call", handler.WebCall)
	// run service
	if err := service.Run(); err != nil {
		log.Error(err.Error())
	}
}

//拦截器
func LanjieqiHandler(e *casbin.Enforcer) gin.HandlerFunc {

	return func(c *gin.Context) {

		//获取请求的URI
		obj := c.Request.URL.RequestURI()
		//获取请求方法
		act := c.Request.Method
		//获取用户的角色
		sub := "admin"

		//判断策略中是否存在
		if e.Enforce(sub, obj, act) {
			log.Info("通过权限")
			c.Next()
		} else {
			log.Info("权限没有通过")
			c.Abort()
		}
	}
}