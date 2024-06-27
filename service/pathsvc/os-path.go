package pathsvc

import (
	"github.com/norniastar/infra-core/contract"
	"path/filepath"
)

type osPath struct {
	root string
}

func (m osPath) GetRoot() string {
	return m.root
}

func (m osPath) Join(paths ...string) string {
	var res string
	for _, path := range paths {
		if res == "" {
			res = path
			continue
		}
		if path == ".." {
			res = filepath.Dir(res)
		} else {
			res = filepath.Join(res, path)
		}
	}
	return res
}

// NewIOPath is 路径实例
func NewIOPath(rootArgs ...string) contract.IIOPath {
	p := new(osPath)
	p.root = p.Join(rootArgs...)
	return p
}
