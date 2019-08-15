package log

type MicroLoggerI struct {}

func (o *MicroLoggerI) Log(v ...interface{}) {
	withCaller(3).Info(v...)
}

// Logf logs formatted using the default logger
func (o *MicroLoggerI) Logf(format string, v ...interface{}) {
	withCaller(3).Infof(format, v...)
}

func NewMicroLogger() *MicroLoggerI {
	return &MicroLoggerI{}
}
