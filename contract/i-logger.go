package contract

// ILogger  Log工厂
type ILogger interface {
	New() ILog
}
