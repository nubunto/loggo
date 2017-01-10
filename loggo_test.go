package loggo

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"testing"
)

func decode(b io.Reader, dst interface{}, unmarshaler func([]byte, interface{}) error) {
	body, _ := ioutil.ReadAll(b)
	unmarshaler(body, dst)
}

func TestNewJSON(t *testing.T) {
	b := new(bytes.Buffer)
	l := New(
		JSON(b),
	)
	l.Msg("hello world!")
	var m map[string]interface{}
	decode(b, &m, json.Unmarshal)
	if m["msg"] != "hello world!" {
		t.Errorf("should be 'hello world!', it is %s", m["msg"])
	}
}
