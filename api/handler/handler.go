package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/client"
	"joe-micro/lib/log"
	"joe-micro/lib/queue"
	"joe-micro/lib/trace"
	srv "joe-micro/service/proto/service"
	"net/http"
	"time"
)

var userClient srv.Service

func init(){
	userClient = srv.NewService("go.micro.srv.service",client.DefaultClient)
}


func WebCall(w http.ResponseWriter, r *http.Request) {
	fmt.Println("调用成功")
	// decode the incoming request as json
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}


	rsp, err := userClient.Call(context.TODO(), &srv.Request{
		Name: request["name"].(string),
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	response := map[string]interface{}{
		"msg": rsp.Msg,
		"ref": time.Now().UnixNano(),
	}


	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}


//测试
func Anything(c *gin.Context){
	c.JSON(200,gin.H{
		"code":0,
		"msg":"调用成功",
	})
}

type   getOne struct {
	Uid   string      `json:"uid" binding:"required"`
}

type User struct {
	Username             string   `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
	Avatar               string   `protobuf:"bytes,2,opt,name=avatar,proto3" json:"avatar,omitempty"`
	Uid                  string   `protobuf:"bytes,3,opt,name=uid,proto3" json:"uid,omitempty"`
	UserType             int32    `protobuf:"varint,4,opt,name=userType,proto3" json:"userType,omitempty"`
}

// @Title Get获取用户信息
// @Description  获取用户信息 Description
// @Accept  json
// @Produce  json
// @Param   token    header    string      true        "token"
// @Param   uid     body    handler.getOne      true        "用户uid"
// @Success 200 {object}  handler.User 	"ok"
// @Router /user/get_one [post]
func GetOne(c *gin.Context){
	log.Info("获取用户信息")
	var req  getOne
	if err := c.ShouldBind(&req); err != nil {
		log.Error(err)
		c.JSON(200,gin.H{
			"code":1,
			"msg":err.Error(),
		})
		return
	}
	ctx, ok := trace.ContextWithSpan(c)
	if !ok {
		log.Warn("不存在context")
		c.JSON(200,gin.H{
			"code":-1,
			"msg":"不存在context",
		})
		return
	}
	rsp, err := userClient.GetOne(ctx, &srv.UserRequest{
		Uid: req.Uid,
	})
	if err != nil {
		log.Error(err)
		c.JSON(200,gin.H{
			"code":-1,
			"msg":err.Error(),
		})
		return
	}
	if  rsp.Uid=="666666"{
		msg := &srv.Message{
			Header: map[string]string{
				"id": fmt.Sprintf("%s", rsp.Uid),
			},
			Body: []byte(fmt.Sprintf("消息 %s ","hello" )),
		}
		body,err:=json.Marshal(&msg)
		if err != nil {
			log.Warn(err)
		}
		if err := queue.Publish("go.micro.srv.service",body); err != nil {
			log.Warn("[pub] 发布消息1失败： %v", err)
		} else {
			log.Info("[pub] 发布消息1：", string(msg.Body))
		}

	}


	c.JSON(200,gin.H{
		"code":0,
		"msg":"",
		"data":rsp,
	})
}

type  cacheM struct {
	Key  string   `json:"key" form:"key"  `
	Value  string    `json:"value" form:"value"`
}

func  PutCache(c *gin.Context){
	log.Info("储存key信息")
	var req  cacheM
	if err := c.ShouldBind(&req); err != nil {
		log.Error(err)
		c.JSON(200,gin.H{
			"code":1,
			"msg":err.Error(),
		})
		return
	}
	rsp, err :=userClient.PutCache(context.TODO(),&srv.CacheRequest{
		Key:req.Key,
		Value:req.Value,
	})
	if err != nil {
		log.Error(err)
		c.JSON(200,gin.H{
			"code":-1,
			"msg":err.Error(),
		})
		return
	}
	c.JSON(200,gin.H{
		"code":0,
		"msg":rsp,
	})
}


func  GetCache(c *gin.Context){
	log.Info("储存key信息")
	var req  cacheM
	if err := c.ShouldBind(&req); err != nil {
		log.Error(err)
		c.JSON(200,gin.H{
			"code":1,
			"msg":err.Error(),
		})
		return
	}
	rsp, err :=userClient.GetCache(context.TODO(),&srv.CacheRequest{
		Key:req.Key,
	})
	if err != nil {
		log.Error(err)
		c.JSON(200,gin.H{
			"code":-1,
			"msg":err.Error(),
		})
		return
	}
	c.JSON(200,gin.H{
		"code":0,
		"msg":"",
		"data":rsp,
	})
}