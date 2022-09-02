package rest

import (
	"io"
	"os"

	"github.com/h14yhv/golang-lib/log"
)

type (
	JsonInterface interface {
		Code(code int) JsonInterface
		Body(data interface{}) JsonInterface
		Log(data interface{}) JsonInterface
		Go() error
	}

	StreamInterface interface {
		Code(code int) StreamInterface
		ContentType(contentType string) StreamInterface
		Body(data io.Reader) StreamInterface
		Go() error
	}

	AttachmentInterface interface {
		Name(name string) AttachmentInterface
		Path(path string) AttachmentInterface
		Go() error
	}
)

var logger log.Logger

func init() {
	logger, _ = log.New(Module, log.DebugLevel, true, os.Stdout)
}
