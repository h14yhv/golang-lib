package rest

import "github.com/labstack/echo/v4"

type attachmentHandler struct {
	ctx  echo.Context
	name string
	path string
}

func Attachment(c echo.Context) AttachmentInterface {
	// Success
	return &attachmentHandler{ctx: c}
}

func (h *attachmentHandler) Name(name string) AttachmentInterface {
	// Success
	h.name = name
	return h
}

func (h *attachmentHandler) Path(path string) AttachmentInterface {
	// Success
	h.path = path
	return h
}

func (h *attachmentHandler) Go() error {
	// Success
	return h.ctx.Attachment(h.path, h.name)
}
