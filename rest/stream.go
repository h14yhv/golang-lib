package rest

import (
	"io"

	"github.com/labstack/echo/v4"
)

type streamHandler struct {
	ctx         echo.Context
	code        int
	contentType string
	body        io.Reader
}

func Stream(c echo.Context) StreamInterface {
	// Success
	return &streamHandler{ctx: c}
}

func (h *streamHandler) Code(code int) StreamInterface {
	// Success
	h.code = code
	return h
}

func (h *streamHandler) ContentType(contentType string) StreamInterface {
	// Success
	h.contentType = contentType
	return h
}

func (h *streamHandler) Body(data io.Reader) StreamInterface {
	// Success
	h.body = data
	return h
}

func (h *streamHandler) Go() error {
	if status := StatusText(h.code); status == "" {
		panic(nil)
	}
	if !ValidContentType(h.contentType) {
		panic(nil)
	}
	// Success
	return h.ctx.Stream(h.code, h.contentType, h.body)
}
