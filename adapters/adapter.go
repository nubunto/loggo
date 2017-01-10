package adapters

import "io"

type Adapter struct {
	Handler
}

func NewAdapter(e Handler) (io.Writer, error) {
	es := &Adapter{
		Handler: e,
	}

	return es, nil
}

func (e *Adapter) Write(b []byte) (int, error) {
	buffercopy := make([]byte, len(b))
	copy(buffercopy, b)
	if err := e.HandlePayload(buffercopy); err != nil {
		return 0, err
	}
	return len(b), nil
}
