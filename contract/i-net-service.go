package contract

type INetService interface {
	Send() (any, error)
	Send_(res any) error
	SetBody(any) INetService
	SetHeaders(map[string]string) INetService
	SetRoute(string, ...any) INetService
	SetURL(string, ...any) INetService
}
