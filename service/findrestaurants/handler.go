package findrestaurants

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
	"restaurants-service/library/errs"
)

type Servicer interface {
	GetListOfRestaurantByKeyword(req Request) (*ResponseModel, error)
}

type Handler struct {
	service Servicer
}

func NewHandler(service Servicer) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) GetListOfRestaurantByKeyword(c echo.Context) error {
	var req Request

	if err := c.Bind(&req); err != nil {
		logrus.Errorf("invalid request: %s", err.Error())
		return errs.JSON(c, errs.New(http.StatusBadRequest, errs.BAD_PARAM.Code, err.Error()))
	}

	if req.Keyword == "" {
		logrus.Error("invalid request: keyword is invalid")
		return errs.JSON(c, errs.New(http.StatusBadRequest, errs.BAD_PARAM.Code, "keyword should not null"))
	}

	res, err := h.service.GetListOfRestaurantByKeyword(req)

	if err != nil {
		logrus.Errorf("call GetListOfRestaurantByKeyword error: %v", err.Error())
		return errs.JSON(c, err)
	}

	return c.JSON(http.StatusOK, res)

}
