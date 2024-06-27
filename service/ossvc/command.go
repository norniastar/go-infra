package ossvc

import (
	"bytes"
	"fmt"
	"os/exec"
	"time"

	"github.com/norniastar/infra-core/contract"
)

var ErrTimeout = fmt.Errorf("执行超时")

type command struct {
	cmd     *exec.Cmd
	expires time.Duration
}

func (m *command) Exec() (string, string, error) {
	var stdout, stderr bytes.Buffer
	m.cmd.Stdout = &stdout
	m.cmd.Stderr = &stderr

	if err := m.cmd.Start(); err != nil {
		return "", "", err
	}

	errChan := make(chan error)
	go func() {
		errChan <- m.cmd.Wait()
	}()

	select {
	case err := <-errChan:
		return stdout.String(), stderr.String(), err
	case <-time.After(m.expires):
		return "", "", ErrTimeout
	}
}

func (m *command) SetDir(format string, args ...interface{}) contract.ICommand {
	m.cmd.Dir = fmt.Sprintf(format, args...)
	return m
}

func (m *command) SetExpires(expires time.Duration) contract.ICommand {
	m.expires = expires
	return m
}
