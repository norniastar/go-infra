package contract

type IEditableAPI interface {
	Call(IUnitOfWork) (any, error)
}
