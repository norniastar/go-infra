package ymlsvc

import (
	json "github.com/bytedance/sonic"
	"github.com/norniastar/go-infra/contract"
)

type configService struct {
	assertDoc map[interface{}]interface{}
	factory   *configFactory
	name      string
}

func (m *configService) Get(key string, value interface{}) error {
	doc := m.getDoc()
	if items, ok := doc[m.name]; ok {
		if item, ok := items.(map[interface{}]interface{})[key]; ok {
			bytes, err := json.Marshal(item)
			if err != nil {
				return err
			}

			return json.Unmarshal(bytes, value)
		}
	}
	return nil
}

func (m *configService) GetStruct(value interface{}) error {
	doc := m.getDoc()
	if temp, ok := doc[m.name]; ok {
		bytes, err := json.Marshal(temp)
		if err != nil {
			return err
		}

		return json.Unmarshal(bytes, value)
	}

	return nil
}

func (m *configService) Has(key string) (ok bool, err error) {
	doc := m.getDoc()
	var items interface{}
	if items, ok = doc[m.name]; ok {
		_, ok = items.(map[interface{}]interface{})[key]
	}

	return
}

func (m *configService) HasStruct() (ok bool, err error) {
	doc := m.getDoc()
	_, ok = doc[m.name]
	return
}

func (m configService) getDoc() map[interface{}]interface{} {
	if m.assertDoc != nil {
		return m.assertDoc
	}

	doc, _ := m.factory.GetDoc()
	return doc
}

func newService(factory *configFactory, name string) contract.IConfigService {
	return &configService{
		factory: factory,
		name:    name,
	}
}
