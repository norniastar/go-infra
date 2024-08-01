package contract

import errorcode "github.com/norniastar/go-infra/model/enum/error-code"

type IError interface {
	error

	GetCode() errorcode.Value
	GetData() any
}
