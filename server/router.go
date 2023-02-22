package server

import (
	"github.com/labstack/echo/v4"
	"restaurants-service/redisclient"
	"restaurants-service/service/findrestaurants"
	"restaurants-service/service/info"
)

func createRoutes(e *echo.Echo, cache *redisclient.Cache) {

	infoHandler := info.NewHandler(
		info.NewService(),
	)

	findRestaurantsService := findrestaurants.NewService(
		redisclient.New(cache.Redis),
	)

	verifyVoucherHandler := findrestaurants.NewHandler(
		findRestaurantsService,
	)

	g := e.Group("/restaurants-service")

	g.GET("/info", infoHandler.Info)

	g.POST("/v1/find", verifyVoucherHandler.FindRestaurant)

}
