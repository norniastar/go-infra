package contract

type IUnitOfWorkRepository interface {
	IUnitOfWork

	RegisterAdd(entry IDbModel)
	RegisterSave(entry IDbModel)
	RegisterRemove(entry IDbModel)
}
