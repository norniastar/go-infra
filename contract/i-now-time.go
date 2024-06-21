package contract

type INowTime interface {
	Unix() int64
	UnixNano() int64
}
