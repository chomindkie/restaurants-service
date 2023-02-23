package server

import (
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"googlemaps.github.io/maps"
	"restaurants-service/redisclient"
	"restaurants-service/service/findrestaurants"
	"restaurants-service/service/info"
)

func createRoutes(e *echo.Echo, cache *redisclient.Cache) {

	var (
		placeApiClient, _ = maps.NewClient(maps.WithAPIKey(viper.GetString("apiKey")))
	)

	infoHandler := info.NewHandler(
		info.NewService(),
	)

	findRestaurantsService := findrestaurants.NewService(
		placeApiClient,
		redisclient.New(cache.Redis),
	)

	verifyVoucherHandler := findrestaurants.NewHandler(
		findRestaurantsService,
	)

	g := e.Group("/restaurants-service")

	g.GET("/info", infoHandler.Info)

	g.POST("/v1/find", verifyVoucherHandler.GetListOfRestaurantByKeyword)

}
