package log

/*
var logger  *logs.BeeLogger

func  init() {
	logger =  logs.NewLogger()
	//输出文件号
	logger.EnableFuncCallDepth(true)
	//异步输出
	logger.Async()
	//异步输出允许设置缓冲 chan 的大小
	//logger.Async(1e3)

	//输出到终端
	logger.SetLogger(logs.AdapterConsole, `{"level":7,"color":true}`)

	//输出到文件
	logger.SetLogger(logs.AdapterFile,`{"filename":"project.log","level":7,"maxlines":0,"maxsize":0,"daily":true,"maxdays":10,"color":true}`)

	//日志直接调用的层级
	logger.SetLogFuncCallDepth(3)

	//输出到 ElasticSearch:
	//logs.SetLogger(logs.AdapterEs, `{"dsn":"http://localhost:9200/","level":1}`)
}*/

/*
级别依次降低，默认全部打印，但是一般我们在部署环境，可以通过设置界别设置日志级别：
LevelEmergency
LevelAlert
LevelCritical
LevelError
LevelWarning
LevelNotice
LevelInformational
LevelDebug*/

/*func Debugf(format string, v ...interface{}) { logger.Debug(format, v) }
func Debug(format string ) { logger.Debug(format) }

func Infof(format string, v ...interface{}) { logger.Info(format, v) }
func Info(format string) { logger.Info(format) }

func Tracef(format string, v ...interface{}) { logger.Trace(format, v) }
func Trace(format string) { logger.Trace(format) }

func Warnf(format string, v ...interface{}) { logger.Warn(format, v) }
func Warn(format string) { logger.Warn(format) }

func Errorf(format string, v ...interface{}) { logger.Error(format, v) }
func Error(format string) { logger.Error(format) }*/

import (
	"fmt"
	"github.com/go-stack/stack"
	"github.com/lestrrat/go-file-rotatelogs"
	"github.com/olivere/elastic"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"gopkg.in/sohlich/elogrus.v3"
	"joe-micro/lib/config"
	"os"
	"time"
)

// 定义键名
const (
	TraceIDKey     = "trace_id"
	UserIDKey      = "user_id"
	VersionKey     = "version"
	ServiceNameKey = "service_name"
)

var logger *logrus.Logger

// debug: 使用text格式, Level是Debug, 打印所有级别
// not debug: 使用json格式, level是Info, 不打印Debug级别
func SetDebug(d bool) {
	if d {
		format := new(logrus.TextFormatter)
		//format.ForceColors = true
		format.TimestampFormat = "2006-01-02 15:04:05"
		logger.Level = logrus.DebugLevel
		logger.Formatter = format
	} else {
		format := new(logrus.JSONFormatter)
		format.TimestampFormat = "2006-01-02 15:04:05"
		logger.Level = logrus.InfoLevel
		logger.Formatter = format
	}
}

func WithField(key string, value interface{}) *logrus.Entry {
	return withCaller(2).WithField(key, value)
}

func WithFields(fs logrus.Fields) *logrus.Entry {
	return withCaller(2).WithFields(fs)
}

//todo 加入traceId 和 uid 日志染色追踪
/*func StartSpan(ctx context.Context, ) *logrus.Entry {
	if ctx == nil {
		ctx = context.Background()
	}
	fields := map[string]interface{}{
		UserIDKey:       fromUserIDContext(ctx),
		TraceIDKey:      fromTraceIDContext(ctx),
		VersionKey:      config.C.Service.Version,
		ServiceNameKey:   config.C.Service.Name,
	}

	return  logger.WithFields(fields)
}


func fromUserIDContext(ctx context.Context) string {
	v := ctx.Value(UserIDKey)
	if v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func fromTraceIDContext(ctx context.Context) string {
	v := ctx.Value(TraceIDKey)
	if v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}*/

func withCaller(skip int) *logrus.Entry {
	var key = "caller"
	var value interface{}
	value = fmt.Sprintf("%+v", stack.Caller(skip))
	return logger.WithFields(logrus.Fields{key: value})
}

/**
 * 使用级别，参照一下
 * - Fatal：网站挂了，或者极度不正常
 * - Error：跟遇到的用户说对不起，可能有bug
 * - Warn：记录一下，某事又发生了
 * - Info：提示一切正常
 * - debug：没问题，就看看堆栈
 **/

func Fatal(args ...interface{}) {
	withCaller(2).Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	withCaller(2).Fatalf(format, args...)
}

func Error(args ...interface{}) {
	withCaller(2).Error(args...)
}

func Errorf(format string, args ...interface{}) {
	withCaller(2).Errorf(format, args...)
}

func Warn(args ...interface{}) {
	withCaller(2).Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	withCaller(2).Warnf(format, args...)
}

func Info(args ...interface{}) {
	withCaller(2).Info(args...)
}

func Infof(format string, args ...interface{}) {
	withCaller(2).Infof(format, args...)
}

func Debug(args ...interface{}) {
	withCaller(2).Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	withCaller(2).Debugf(format, args...)
}

// 输出日志到es
func configESLogger(esUrl string, esHOst string, index string) {
	client, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(esUrl))
	if err != nil {
		logger.Errorf("config es logger error. %+v", errors.WithStack(err))
		return
	}
	esHook, err := elogrus.NewElasticHook(client, esHOst, logrus.DebugLevel, index)
	if err != nil {
		logger.Errorf("config es logger error. %+v", errors.WithStack(err))
		return
	}

	logger.AddHook(esHook)
}

//输出日志到文件
func configFileLogger(logPrefix string) {
	logWriter, _ := rotatelogs.New(
		logPrefix+"_%Y-%m-%d.log",
		rotatelogs.WithMaxAge(7*24*time.Hour),     // 文件最大保存时间
		rotatelogs.WithRotationTime(24*time.Hour), // 日志切割时间间隔
	)
	writeMap := lfshook.WriterMap{
		logrus.InfoLevel: logWriter,
		//logrus.FatalLevel: logWriter, //错误输出到另一个日志
	}
	lfHook := lfshook.NewHook(writeMap, &logrus.JSONFormatter{TimestampFormat: "2006-01-02 15:04:05"})

	logger.AddHook(lfHook)
}

func init() {
	logger = &logrus.Logger{
		Out:       os.Stdout,
		Formatter: nil,
		Hooks:     make(logrus.LevelHooks),
		Level:     0,
	}
	SetDebug(config.C.Debug)
	//ConfigESLogger("http://localhost:9200","localhost","mylog")
	//configFileLogger("api")
}
