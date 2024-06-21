package contract

type IResponseWriter interface {
	Write(writer any, data any)
}
