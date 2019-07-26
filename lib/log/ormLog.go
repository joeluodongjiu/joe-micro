package log

type OrmLoggerI struct {
}

func (o *OrmLoggerI) Print(v ...interface{}) {
	format := v[0].(string)
	v = v[2:]  //只取sql语句
	//跳8级
	withCaller(8).Infof(format+" %v  ", v)
}

func NewGormLogger() *OrmLoggerI {
	return &OrmLoggerI{}
}
