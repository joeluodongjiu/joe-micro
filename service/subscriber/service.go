package subscriber

import (
	"encoding/json"
	"github.com/nsqio/go-nsq"
	"joe-micro/lib/log"
	"joe-micro/lib/queue"
	service "joe-micro/service/proto/service"
)

func Registersubscriber() {
	//可监听多个主题
	err := queue.Subscribe("go.micro.srv.service", "sayHello", new(SayHello))
	if err != nil {
		log.Fatal(err)
		return
	}

/*	err = queue.Subscribe("go.micro.srv.service2", "sayHello", new(SayHello))
	if err != nil {
		log.Fatal(err)
		return
	}*/
}


type SayHello struct{}

func (s SayHello) HandleMessage(m *nsq.Message) error {
	log.Info("读取队列")
	//return errors.New("错误")
	msg := service.Message{}
	err := json.Unmarshal(m.Body, &msg)
	if err != nil {
		log.Warn(err)
		return err
	}
	log.Infof("header: %v   body: %v ", msg.Header, string(msg.Body))
	return nil
}
