package ossvc

import (
	"os"
	"path/filepath"

	"github.com/norniastar/infra-core/contract"
)

type ioNode struct {
	ioPath contract.IIOPath
	path   string
}

func (m ioNode) GetName() string {
	return filepath.Base(m.path)
}

func (m ioNode) GetParent() contract.IIODirectory {
	return NewIODirectory(
		m.ioPath,
		m.GetPath(),
		"..",
	)
}

func (m ioNode) GetPath() string {
	return m.path
}

func (m ioNode) IsExist() bool {
	_, err := os.Stat(m.path)
	return err == nil || os.IsExist(err)
}

func (m ioNode) IsDir() bool {
	p, err := os.Stat(m.path)
	return err == nil && p.IsDir()
}

func (m ioNode) Move(paths ...string) error {
	dstPath := m.ioPath.Join(paths...)
	return os.Rename(
		m.GetPath(),
		dstPath,
	)
}

func (m ioNode) Remove() error {
	if !m.IsExist() {
		return nil
	}

	return os.RemoveAll(
		m.GetPath(),
	)
}

func newIONode(ioPath contract.IIOPath, paths ...string) contract.IIONode {
	return ioNode{
		ioPath: ioPath,
		path:   ioPath.Join(paths...),
	}
}
