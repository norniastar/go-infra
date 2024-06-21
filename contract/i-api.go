package contract

type IAPI interface {
	Call() (any, error)
}
