package log

import (
	"strings"
)

// nsq的专属log
// 1. 不关心caller
// 2. level需要处理

type NsqLoggerI struct {}

func (l *NsqLoggerI) Output(calldepth int, s string) error {
	level := strings.Split(s, " ")[0]
	switch level {
	case "INF":
		withCaller(2).Info(s)
	case "WRN":
		withCaller(2).Warn(s)
	case "ERR":
		withCaller(2).Error(s)
	default:
		withCaller(2).Debug(s)
	}
	return nil
}

func NsqLogger() *NsqLoggerI {
	return &NsqLoggerI{}
}
