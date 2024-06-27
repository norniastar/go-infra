package contract

type IIONode interface {
	GetName() string
	GetParent() IIODirectory
	GetPath() string
	IsExist() bool
	IsDir() bool
	Move(pathArgs ...string) error
	Remove() error
}
