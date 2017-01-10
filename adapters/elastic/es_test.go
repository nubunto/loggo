package elastic

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"testing"

	"github.com/nubunto/loggo"
	"github.com/nubunto/loggo/adapters"
)

func TestESAdapter(t *testing.T) {
	esBuffer := new(bytes.Buffer)
	esHandler := adapters.HandlerFunc(func(b []byte) error {
		esBuffer.Write(b)
		return nil
	})

	adapter, err := adapters.NewAdapter(esHandler)
	if err != nil {
		t.Fatal(err)
	}
	l := loggo.New(loggo.JSON(adapter))

	if err := l.Log("this-should", "go to esBuffer"); err != nil {
		t.Fatalf("Shouldn't error: instead, we have %#v", err)
	}
	m := make(map[string]interface{})
	body, _ := ioutil.ReadAll(esBuffer)
	json.Unmarshal(body, &m)
	if m["this-should"] != "go to esBuffer" {
		t.Errorf("m['this-should'] should be 'go to esBuffer', have %s instead", m["this-should"])
	}
}

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

func TestElasticAdapterErrHandler(t *testing.T) {
	sentinel := errors.New("oops")
	esHandler := adapters.HandlerFunc(func(b []byte) error {
		return sentinel
	})

	esAdapter, err := adapters.NewAdapter(esHandler)
	if err != nil {
		t.Fatal(err)
	}
	l := loggo.New(
		loggo.JSON(esAdapter),
	)
	err = l.Log("this", "should fail")
	if err != sentinel {
		t.Errorf("error %v should be %v", err, sentinel)
	}
}
