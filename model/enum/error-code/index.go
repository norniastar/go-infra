package errorcode

type Value int

const (
	Null    Value = 0   // Null  无效错误
	API     Value = 501 // API  api错误码
	Auth    Value = 502 // Auth  认证错误码
	Verify  Value = 503 // Verify  验证错误码
	Tip     Value = 504 // Tip  提醒错误码
	Timeout Value = 505 // Timeout  超时
	Panic   Value = 599 // Panic  异常错误码
)
