package main

import (
	"github.com/micro/cli"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/consul"
	microLog "github.com/micro/go-micro/util/log"
	"github.com/micro/go-micro/web"
	"github.com/opentracing/opentracing-go"
	"joe-micro/adminApi/routers"
	"joe-micro/lib/config"
	"joe-micro/lib/log"
	"joe-micro/lib/queue"
	"joe-micro/lib/trace"
	"time"
)

// @title  微服务的管理端api文档demo
// @version 1.0
// @host  localhost:9081
// @BasePath /api/admin

// @securityDefinitions.apikey MustToken
// @in header
// @name token


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


	/************************************/
	/********** gin  路由框架     ********/
	/************************************/
	//注册 gin  routers
	ginHandler :=routers.Init()


	// create new api service
	service := web.NewService(
		web.Name(config.C.Service.Name),
		web.Registry(reg),
		web.RegisterTTL(time.Second*15),      //重新注册时间
		web.RegisterInterval(time.Second*10), //注册过期时间
		web.Version(config.C.Service.Version),
		web.Address(config.C.Service.Port),
		web.Handler(ginHandler),
	)

	// initialise service
	if err := service.Init(
		web.Action(func(ctx *cli.Context) {
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

		}),
	); err != nil {
		log.Error(err.Error())
	}



	// run service
	if err := service.Run(); err != nil {
		log.Error(err.Error())
	}
}
