package main

import "C"
import (
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/consul"
	"github.com/micro/go-micro/web"
	"github.com/opentracing/opentracing-go"
	_ "joe-micro/api/docs"
	"joe-micro/api/handler"
	"joe-micro/lib/config"
	"joe-micro/lib/log"
	"joe-micro/lib/queue"
	"joe-micro/lib/trace"
)

// @title  微服务的api文档demo
// @version 1.0
// @host  localhost:8081
// @BasePath /
func main() {

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
	queue.Init(config.C.Nsq.Address,config.C.Nsq.Lookup,config.C.Nsq.MaxInFlight,config.C.Debug)



	/************************************/
	/********** gin  路由框架     ********/
	/************************************/
	gin.SetMode(gin.ReleaseMode)   //是否生产模式启动
	router:=gin.Default()
	router.Use(log.GinLogger())
	router.Use(trace.TracerWrapper)
	// 跨域
	router.Use(func(ctx *gin.Context) {
		ctx.Header("Access-Control-Allow-Origin", "*")   //跨域
		ctx.Header("Access-Control-Allow-Headers", "Token,Content-Type")  //必须的请求头
		ctx.Header("Access-Control-Allow-Methods", "OPTIONS,PUT,POST,GET,DELETE")  //接收的请求方法
	})

	// OPTIONS
	router.NoRoute(func(ctx *gin.Context) {
		if ctx.Request.Method == "OPTIONS" {
			ctx.JSON(200, nil)
		}
	})

	//swagger
/*	url := ginSwagger.URL("http://localhost:8081/swagger/doc.json") // The url pointing to API definition
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))*/

	r := router.Group("/user")
	r.GET("/test", handler.Anything)

	r.POST("/get_one", handler.GetOne)
	r.POST("/put_cache",handler.PutCache)
	r.GET("/get_cache",handler.GetCache)

	// register html handler
	service.Handle("/", router)

	// register call handler

	service.HandleFunc("/api/call", handler.WebCall)
	// run service
	if err := service.Run(); err != nil {
		log.Error(err.Error())
	}
}
