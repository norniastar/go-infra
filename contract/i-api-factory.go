package contract

type IAPIFactory interface {
	Build(endpoint, name, version string) any
	Register(endpoint, name string, api any)
}
