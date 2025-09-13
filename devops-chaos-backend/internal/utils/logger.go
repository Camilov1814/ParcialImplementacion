// internal/utils/logger.go
package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Logger = logrus.New()

func InitLogger() {
	Logger.SetOutput(os.Stdout)
	Logger.SetLevel(logrus.InfoLevel)
	if os.Getenv("GIN_MODE") == "debug" {
		Logger.SetLevel(logrus.DebugLevel)
	}
}
