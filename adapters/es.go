package adapters

import (
	"io"

	elastic "gopkg.in/olivere/elastic.v3"
)

const (
	defaultIndex = "elastic-logger"
	defaultType  = "elastic-log"
)

type Adapter struct {
	Handler

	errHandler func(error)
}

type AdapterOption func(*Adapter) error

type ElasticHandler struct {
	index    string
	typeName string
	client   *elastic.Client
}

type ElasticHandlerOption func(*ElasticHandler) error

func Index(index string) ElasticHandlerOption {
	return func(e *ElasticHandler) error {
		e.index = index
		return nil
	}
}

func Type(t string) ElasticHandlerOption {
	return func(e *ElasticHandler) error {
		e.typeName = t
		return nil
	}
}

func Client(opts ...elastic.ClientOptionFunc) ElasticHandlerOption {
	return func(e *ElasticHandler) error {
		c, err := elastic.NewClient(opts...)
		if err != nil {
			return err
		}
		e.client = c
		return nil
	}
}

func (e *ElasticHandler) HandlePayload(b []byte) error {
	_, err := e.client.Index().
		Index(e.index).
		Type(e.typeName).
		BodyString(string(b)).
		Do()
	if err != nil {
		return err
	}
	return nil
}

func NewElasticHandler(opts ...ElasticHandlerOption) (*ElasticHandler, error) {
	e := &ElasticHandler{
		index:    defaultIndex,
		typeName: defaultType,
	}

	for _, opt := range opts {
		if err := opt(e); err != nil {
			return nil, err
		}
	}
	return e, nil
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
