package errorsvc

import (
	"fmt"
	"github.com/norniastar/infra-core/contract"
	errorcode "github.com/norniastar/infra-core/model/enum/error-code"
)

type custom struct {
	error

	code errorcode.Value
	data any
}

func (m custom) Error() string {
	return fmt.Sprintf("[err: %v, code: %v, data: %v]", m.error, m.code, m.data)
}

func (m custom) GetCode() errorcode.Value {
	return m.code
}

func (m custom) GetData() any {
	return m.data
}

func New(code errorcode.Value, data any) contract.IError {
	return custom{
		error: fmt.Errorf("%v", data),
		code:  code,
		data:  data,
	}
}

func Newf(code errorcode.Value, format string, args ...any) contract.IError {
	return custom{
		error: fmt.Errorf(format, args),
		code:  code,
		data:  fmt.Errorf(format, args),
	}
}

func NewError(code errorcode.Value, err error) contract.IError {
	return custom{
		error: err,
		code:  code,
	}
}

func NewErrorf(code errorcode.Value, format string, args ...any) contract.IError {
	return custom{
		error: fmt.Errorf(format, args),
		code:  code,
	}
}
