package server

import (
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	logger    = logrus.StandardLogger()
	FILE_PATH = os.Getenv("CONFIG_PATH")
)

func InitViper(defaultFilePath string) {
	if FILE_PATH == "" {
		FILE_PATH = defaultFilePath
	}

	viper.AddConfigPath(getAbsFilePath(FILE_PATH))
	viper.SetConfigName("application")
	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		logger.Errorf("Cannot read config: %s", err)
	}
	viper.AutomaticEnv()
}

func getAbsFilePath(path string) string {
	res, err := filepath.Abs(path)
	if err != nil {
		return path
	}
	return res
}
