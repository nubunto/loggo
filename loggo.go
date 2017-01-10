package loggo

import (
	"io"

	"github.com/go-kit/kit/log"
)

type Logger struct {
	log.Logger
}

func (l Logger) Msg(v string) error {
	return l.Log("msg", v)
}

func New(underlyingLogger log.Logger) Logger {
	return Logger{
		underlyingLogger,
	}
}

func JSON(output io.Writer) log.Logger {
	return log.NewJSONLogger(output)
}
