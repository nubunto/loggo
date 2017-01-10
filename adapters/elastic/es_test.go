package elastic

import (
	"testing"

	"github.com/nubunto/loggo"
	"github.com/nubunto/loggo/adapters"
)

func TestElasticHandlerOptions(t *testing.T) {
	esHandler, err := NewElasticHandler(Index("foobar"), Type("not-a-log"))
	if err != nil {
		t.Fatal(err)
	}
	c := esHandler.(*ElasticHandler)
	if c.index != "foobar" {
		t.Errorf("index should be foobar, it is %s", c.index)
	}
	if c.typeName != "not-a-log" {
		t.Errorf("typeName should be 'not-a-log', it is %s", c.typeName)
	}
}

type errLogger struct {
	err error
	loggo.Logger
}

func (e errLogger) Log(args ...interface{}) error {
	if e.err != nil {
		// stop
		return e.err
	}
	if err := e.Logger.Log(args...); err != nil {
		e.err = err
		return err
	}
	return nil
}

func TestRealES(t *testing.T) {
	esHandler, err := NewElasticHandler(Client())
	if err != nil {
		t.Fatal(err)
	}
	esAdapter, err := adapters.NewAdapter(esHandler)
	if err != nil {
		t.Fatal(err)
	}
	esLogger := loggo.New(
		loggo.JSON(esAdapter),
	)
	e := errLogger{Logger: esLogger}
	e.Log("a", 1)
	e.Log("b", 2)
	e.Log("c", 3)
	if e.err != nil {
		t.Fatal(err)
	}
}
