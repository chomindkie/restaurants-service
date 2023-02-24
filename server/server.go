package server

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4/middleware"
	"os"
	"os/signal"
	"restaurants-service/library/errs"
	"restaurants-service/redisclient"
	"syscall"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

func Start() {
	var (
		e     = initEcho()
		cache = newCache()
	)

	defer func() {
		cache.Redis.Close()
	}()

	createRoutes(e, cache)

	go func() {
		servicePort := viper.GetString("service.port")
		serviceCert := viper.GetString("service.cert-file")
		serviceKey := viper.GetString("service.key-file")
		fmt.Printf("Starting application server on port: %s", servicePort)
		e.Logger.Fatal(e.StartTLS(":"+servicePort, getAbsFilePath(serviceCert), getAbsFilePath(serviceKey)))
	}()

	gracefulShutdown(e)
}

func initEcho() *echo.Echo {
	// Echo instance
	e := echo.New()
	e.HideBanner = true
	e.HTTPErrorHandler = errs.HTTPErrorHandler

	// Middleware
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	return e
}

func newCache() *redisclient.Cache {
	return redisclient.NewCache(viper.GetString("redis.host"), viper.GetInt("redis.db-index"))
}

func gracefulShutdown(e *echo.Echo) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	if err := e.Shutdown(context.Background()); err != nil {
		logger.Printf("shutdown server: %s", err.Error())
	}
}
