package handler

import (
	"context"
	"joe-micro/lib/cache"
	"joe-micro/lib/log"
	"joe-micro/service/model"
	service "joe-micro/service/proto/service"
	"time"
)

type Service struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Service) Call(ctx context.Context, req *service.Request, rsp *service.Response) error {
	log.Info("Received Service.Call request")
	rsp.Msg = "Hello " + req.Name
	return nil
}


// Call is a single request handler called via client.Call or the generated client code
func (e *Service) GetOne(ctx context.Context, req *service.UserRequest, rsp *service.User) (err error) {
	// simulate opentracing instrumentation of an SQL query

	log.Info("Received Service.GetOne request",)
	var user model.User
	user,err = model.GetByID(req.Uid)
	if err!=nil {
		log.Warn(err.Error())
		return err
	}
	rsp.Uid = user.Uid
	rsp.Username = user.Username
	rsp.Avatar = user.Avatar
	rsp.UserType = int32(user.UserType)
	return nil
}


// Call is a single request handler called via client.Call or the generated client code
func (e *Service) PutCache(ctx context.Context, req *service.CacheRequest, rsp *service.CacheResponse) error {
	log.Info("Received Service.PutCache request")
	err:=cache.Cache.Put(req.Key,req.Value,50*time.Second)
	if err!=nil {
		log.Error(err)
		return err
	}
	rsp.Value="成功"
	return nil
}

// Call is a single request handler called via client.Call or the generated client code
func (e *Service) GetCache(ctx context.Context, req *service.CacheRequest, rsp *service.CacheResponse) error {
	rsp.Value =cache.GetString(cache.Cache.Get(req.Key))
	return nil
}





// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *Service) Stream(ctx context.Context, req *service.StreamingRequest, stream service.Service_StreamStream) error {
	log.Infof("Received Service.Stream request with count: %d", req.Count)

	for i := 0; i < int(req.Count); i++ {
		log.Warnf("Responding: %d", i)
		if err := stream.Send(&service.StreamingResponse{
			Count: int64(i),
		}); err != nil {
			return err
		}
	}

	return nil
}

// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
func (e *Service) PingPong(ctx context.Context, stream service.Service_PingPongStream) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Infof("Got ping %v", req.Stroke)
		if err := stream.Send(&service.Pong{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}
