package errorsvc

import (
	"fmt"
	"github.com/norniastar/infra-core/contract"
	"runtime/debug"
)

func Recover(log contract.ILog, err *error, handleErrorAction func(error)) {
	var cErr error
	if rv := recover(); rv != nil {
		var ok bool
		if cErr, ok = rv.(error); !ok {
			cErr = fmt.Errorf("%v", rv)
		}

		if log != nil {
			log.AddLabel(
				"stack",
				"%s",
				debug.Stack(),
			)
		}
	}

	if cErr == nil && err != nil {
		cErr = *err
	}

	if cErr != nil {
		handleErrorAction(cErr)
	}
}
