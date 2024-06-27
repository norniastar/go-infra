package contract

import "reflect"

// ITraceable 跟踪接口
type ITraceable interface {
	WithContext(ctx any) reflect.Value
}
