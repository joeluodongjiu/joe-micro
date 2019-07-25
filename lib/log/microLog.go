package log

import "github.com/sirupsen/logrus"

type MicroLoggerI struct {
	logger *logrus.Logger
}

func (o *MicroLoggerI) Log(v ...interface{}) {
	o.logger.Info(v)
}

// Logf logs formatted using the default logger
func (o *MicroLoggerI) Logf(format string, v ...interface{}) {
	o.logger.Infof(format, v...)
}

func NewMicroLogger() *MicroLoggerI {
	return &MicroLoggerI{logger: logger}
}
