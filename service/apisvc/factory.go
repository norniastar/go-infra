package apisvc

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/norniastar/go-infra/contract"
	errorcode "github.com/norniastar/go-infra/model/enum/error-code"
	"github.com/norniastar/go-infra/service/errorsvc"
)

var (
	errNilApi = errorsvc.Newf(errorcode.API, "")
	nilApiPtr = &nilApi{}
)

type factory map[string]map[string]reflect.Type

func (m factory) Build(endpoint, name, version string) any {
	if apiTypes, ok := m[endpoint]; ok {
		names := [2]string{"", name}
		if version != "" {
			names[0] = fmt.Sprintf(
				"%s_%s",
				name,
				strings.Replace(version, ".", "_", -1),
			)
		}

		for _, r := range names {
			if apiType, ok := apiTypes[r]; ok {
				return reflect.New(apiType).Interface()
			}
		}
	}

	return nilApiPtr
}

func (m factory) Register(endpoint, name string, api any) {
	if _, ok := m[endpoint]; !ok {
		m[endpoint] = make(map[string]reflect.Type)
	}

	apiType := reflect.TypeOf(api)
	if apiType.Kind() == reflect.Ptr {
		apiType = apiType.Elem()
	}
	m[endpoint][name] = apiType
}

type nilApi struct{}

func (m nilApi) Call() (any, error) {
	return nil, errNilApi
}

func NewFactory() contract.IAPIFactory {
	return make(factory)
}
