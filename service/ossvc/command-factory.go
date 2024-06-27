package ossvc

import (
	"os/exec"
	"time"

	"github.com/norniastar/infra-core/contract"
)

type factory func(name string, args ...string) contract.ICommand

func (m factory) Build(name string, args ...string) contract.ICommand {
	return m(name, args...)
}

func NewCommandFactory() contract.ICommandFactory {
	return factory(func(name string, args ...string) contract.ICommand {
		return &command{
			cmd:     exec.Command(name, args...),
			expires: 4 * time.Second,
		}
	})
}
