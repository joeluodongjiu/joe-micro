package main

import (
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/consul"
	ocplugin "github.com/micro/go-plugins/wrapper/trace/opentracing"
	"github.com/opentracing/opentracing-go"
	"joe-micro/lib/config"
	"joe-micro/lib/log"
	"joe-micro/lib/queue"
	"joe-micro/lib/trace"
	"joe-micro/service/handler"
	srv "joe-micro/service/proto/service"
	"joe-micro/service/subscriber"
	"time"
)

func main() {

	/************************************/
	/********** 服务发现  cousul   ********/
	/************************************/
	reg := consul.NewRegistry(func(op *registry.Options) {
		op.Addrs = []string{
			config.C.Consul,
		}
		op.Timeout =  5 * time.Second
	})




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
	subscriber.Registersubscriber()  //注册消费者


	// New Service
	service := micro.NewService(
		micro.Name(config.C.Service.Name),
		micro.Registry(reg),
		micro.RegisterTTL(time.Second*15),        //重新注册时间
		micro.RegisterInterval(time.Second*10),   //注册过期时间
		micro.Version(config.C.Service.Version),
		micro.WrapHandler(ocplugin.NewHandlerWrapper(opentracing.GlobalTracer())),

	)

	// Initialise service
	service.Init()

	// Register Handler
	err=srv.RegisterServiceHandler(service.Server(), new(handler.Service))
	if err != nil {
		log.Fatal(err)
	}




	// Run service

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
