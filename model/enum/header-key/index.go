package headerkey

type Value string

func (m Value) String() string {
	return string(m)
}

const (
	AccessToken   Value = "H-A"           // AccessToken  访问令牌
	Debug         Value = "H-DEBUG"       // Debug  调试
	Endpoint      Value = "H-E"           // Endpoint  端
	MobileSession Value = "H-M-S"         // MobileSession  移动会话
	Timeout       Value = "H-O"           // Timeout  超时
	TraceID       Value = "Uber-Trace-Id" // TraceID  追踪头
	Version       Value = "H-V"           // Version  版本
)
