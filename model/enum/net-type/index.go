package nettype

type Value int

const (
	HTTPPost Value = 1 // HTTPPost is http post
	MQ       Value = 2 // MQ is 消息队列
	GRPC     Value = 3
)
