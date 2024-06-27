package yamlsvc

import (
	"reflect"
	"sync"

	"github.com/norniastar/infra-core/contract"
	configgroup "github.com/norniastar/infra-core/model/enum/config-group"
)

var configFactoryMutex sync.Mutex

type configFactory struct {
	doc  map[interface{}]interface{}
	file contract.IIOFile
}

func (m *configFactory) Build(group configgroup.Value) contract.IConfigService {
	return newService(m, group.String())
}

func (m *configFactory) Build_(s interface{}) contract.IConfigService {
	return newService(
		m,
		reflect.TypeOf(s).Name(),
	)
}

func (m *configFactory) GetDoc() (map[interface{}]interface{}, error) {
	if m.doc == nil {
		configFactoryMutex.Lock()
		defer configFactoryMutex.Unlock()

		if m.doc == nil {
			if err := m.file.ReadYaml(&(m.doc)); err != nil {
				return nil, err
			}
		}
	}

	return m.doc, nil
}

func NewConfigFactory(file contract.IIOFile) contract.IConfigFactory {
	return &configFactory{
		file: file,
	}
}
