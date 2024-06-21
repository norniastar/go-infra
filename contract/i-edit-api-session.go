package contract

type IEditableAPISession interface {
	SetSession(IUnitOfWork, any) error
}
