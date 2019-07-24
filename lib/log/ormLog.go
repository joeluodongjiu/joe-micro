package log

import (
	"github.com/sirupsen/logrus"
)

type  OrmLoggerI struct {
	logger    *logrus.Logger
}



func (o *OrmLoggerI)Print(v ...interface{})  {
	format := v[0].(string)
	v = v[1:]
	o.logger.Infof(format+" %v  ",v)
}

func NewGormLogger()  *OrmLoggerI {
	return   &OrmLoggerI{logger:logger}
}

