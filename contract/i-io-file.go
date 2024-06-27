package contract

import "os"

type IIOFile interface {
	IIONode

	GetExt() string
	GetFile() (*os.File, error)
	Read(data interface{}) error
	ReadJSON(data interface{}) error
	ReadYaml(data interface{}) error
	Write(data interface{}) error
}
