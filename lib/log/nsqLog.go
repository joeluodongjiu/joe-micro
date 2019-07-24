package log

import (
	"github.com/sirupsen/logrus"
	"strings"
)

// nsq的专属log
// 1. 不关心caller
// 2. level需要处理

type NsqLoggerI struct {
	logger *logrus.Logger
}

func (l *NsqLoggerI) Output(calldepth int, s string) error {
	level := strings.Split(s, " ")[0]
	switch level {
	case "INF":
		l.logger.Info(s)
	case "ERR":
		l.logger.Error(s)
	default:
		l.logger.Info(s)
	}
	return nil
}

func NsqLogger() *NsqLoggerI {
	return &NsqLoggerI{logger: logger}
}
