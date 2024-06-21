package contract

// IUnitOfWork  工作单元
type IUnitOfWork interface {
	Commit() error
}
