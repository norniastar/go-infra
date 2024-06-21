package contract

// ILog  日志类
type ILog interface {
	AddLabel(key, format string, v ...any) ILog
	Debug()
	Error(err error)
	Fatal()
	Info()
	Warning()
}
