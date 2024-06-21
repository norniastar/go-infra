package contract

import "time"

// ILock  锁接口
type ILock interface {
	Lock(key string, expires time.Duration) (func() error, error)
}
