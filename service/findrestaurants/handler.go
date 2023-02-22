package findrestaurants

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"googlemaps.github.io/maps"
	"net/http"
	"restaurants-service/library/errs"
)

type Servicer interface {
	FindRestaurant(c echo.Context, req Request) (*maps.PlacesSearchResponse, error)
}

type Handler struct {
	service Servicer
}

func NewHandler(service Servicer) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) FindRestaurant(c echo.Context) error {
	var req Request

	if err := c.Bind(&req); err != nil {
		log.Errorf("invalid request: %s", err.Error())
		return errs.JSON(c, errs.New(http.StatusBadRequest, errs.BAD_PARAM.Code, err.Error()))
	}

	if req.Keyword == "" {
		log.Error("invalid request: Mobile is invalid")
		return errs.JSON(c, errs.New(http.StatusBadRequest, errs.BAD_PARAM.Code, "keyword should not null"))
	}

	log.Info("Find Restaurant at ", req.Keyword)
	res, err := h.service.FindRestaurant(c, req)

	if err != nil {
		log.Errorf("call FindRestaurant error: ", err.Error())
		return errs.JSON(c, err)
	}

	return c.JSON(http.StatusOK, res)

}
