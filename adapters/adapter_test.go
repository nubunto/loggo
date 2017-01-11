package adapters

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"testing"

	"github.com/nubunto/loggo"
)

func TestESAdapter(t *testing.T) {
	buffer := new(bytes.Buffer)
	handler := HandlerFunc(func(b []byte) error {
		buffer.Write(b)
		return nil
	})

	adapter := NewAdapter(handler)
	l := loggo.New(loggo.JSON(adapter))

	if err := l.Log("this-should", "go to esBuffer"); err != nil {
		t.Fatalf("Shouldn't error: instead, we have %#v", err)
	}
	m := make(map[string]interface{})
	body, _ := ioutil.ReadAll(buffer)
	json.Unmarshal(body, &m)
	if m["this-should"] != "go to esBuffer" {
		t.Errorf("m['this-should'] should be 'go to esBuffer', have %s instead", m["this-should"])
	}
}

func TestElasticAdapterErrHandler(t *testing.T) {
	sentinel := errors.New("oops")
	handler := HandlerFunc(func(b []byte) error {
		return sentinel
	})

	adapter := NewAdapter(handler)
	l := loggo.New(
		loggo.JSON(adapter),
	)
	err := l.Log("this", "should fail")
	if err != sentinel {
		t.Errorf("error %v should be %v", err, sentinel)
	}
}
