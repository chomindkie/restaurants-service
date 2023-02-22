package info

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

var appVersion string = os.Getenv("APP_VERSION")

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Info(c echo.Context) *InfoResponse {
	if appVersion == "" {
		appVersion = viper.GetString("app.version")
	}

	return &InfoResponse{
		Build: Build{
			Name:    viper.GetString("app.name"),
			Version: appVersion,
		},
	}
}
