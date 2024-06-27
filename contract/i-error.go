package contract

import errorcode "github.com/norniastar/infra-core/model/enum/error-code"

type IError interface {
	error

	GetCode() errorcode.Value
	GetData() any
}
