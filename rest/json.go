package rest

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type jsonHandler struct {
	ctx  echo.Context
	code int
	body interface{}
}

var threshold = http.StatusBadRequest

func JSON(c echo.Context) JsonInterface {
	// Success
	return &jsonHandler{ctx: c}
}

func (h *jsonHandler) Code(code int) JsonInterface {
	// Success
	h.code = code
	return h
}

func (h *jsonHandler) Body(data interface{}) JsonInterface {
	// Success
	h.body = data
	return h
}

func (h *jsonHandler) Log(data interface{}) JsonInterface {
	if h.code < threshold {
		logger.Infof("code %d: %v", h.code, data)
	} else {
		logger.Errorf("code %d: %v", h.code, data)
	}
	// Success
	return h
}

func (h *jsonHandler) Go() error {
	status := StatusText(h.code)
	if status == "" {
		panic(nil)
	}
	// Success
	return h.ctx.JSON(h.code, &response{
		Success: Success(h.code),
		Message: status,
		Detail:  h.body,
	})
}
