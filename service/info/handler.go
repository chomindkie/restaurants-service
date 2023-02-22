package info

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Servicer interface {
	Info(c echo.Context) *InfoResponse
}

type Handler struct {
	service Servicer
}

func NewHandler(service Servicer) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Info(c echo.Context) error {
	res := h.service.Info(c)
	return c.JSON(http.StatusOK, res)
}
