package contract

type IUnitOfWorkEvent interface {
	RegisterAfterCommit(func())
}
