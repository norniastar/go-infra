package contract

import "time"

// ICommand is 命令接口
type ICommand interface {
	Exec() (stdout string, stderr string, err error)
	SetDir(format string, args ...interface{}) ICommand
	SetExpires(expires time.Duration) ICommand
}
