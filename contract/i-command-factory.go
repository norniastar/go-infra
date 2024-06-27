package contract

type ICommandFactory interface {
	Build(name string, args ...string) ICommand
}
