package log

import "github.com/sirupsen/logrus"

type OrmLoggerI struct {
}

func (o *OrmLoggerI) Print(v ...interface{}) {
	format := v[0].(string)
	caller:= v[1].(string)
	v = v[2:]  //只取sql语句
	//跳8级
	logger.WithFields(logrus.Fields{"caller": caller}).Infof(format+" %v  ", v)
}

func NewGormLogger() *OrmLoggerI {
	return &OrmLoggerI{}
}
