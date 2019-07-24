package subscriber

import (
	"encoding/json"
	"github.com/nsqio/go-nsq"
	"joe-micro/lib/log"
	"joe-micro/lib/queue"
	service "joe-micro/service/proto/service"
)

func Registersubscriber() {
	sub1("go.micro.srv.service")

}

func sub1(topic string) {
	err := queue.Subscribe(topic,"sayHello",new(SayHello))
	if err != nil {
		log.Fatal(err)
		return
	}
}



type   SayHello struct {}

func (s SayHello)HandleMessage(m *nsq.Message) error{
	msg := service.Message{}
	err:=json.Unmarshal(m.Body,&msg)
	if err != nil {
		log.Warn(err)
		return err
	}
	log.Infof("header: %v   body: %v ",msg.Header,string(msg.Body))
	return nil
}