package contract

type IIODirectory interface {
	IIONode

	Create() error
	FindDirectories() []IIODirectory
	FindFiles() []IIOFile
}
