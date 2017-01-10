package adapters

type Handler interface {
	HandlePayload([]byte) error
}

type HandlerFunc func([]byte) error

func (e HandlerFunc) HandlePayload(b []byte) error {
	return e(b)
}
