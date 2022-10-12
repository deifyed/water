package logging

import (
	"fmt"
	"os"

	"github.com/deifyed/water/pkg/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func GetLogger() *logrus.Logger {
	log := logrus.New()

	var err error

	log.Out = os.Stdout
	log.Formatter = &logrus.JSONFormatter{PrettyPrint: true}

	log.Level, err = parseLevel(viper.GetString(config.LogLevel))
	if err != nil {
		panic("test")
	}

	return log
}

func parseLevel(level string) (logrus.Level, error) {
	fmt.Println(level)
	fmt.Println(viper.GetString(config.LogLevel))
	fmt.Println(viper.ConfigFileUsed())
	fmt.Println(viper.AllKeys())

	switch level {
	case "debug":
		return logrus.DebugLevel, nil
	case "info":
		return logrus.InfoLevel, nil
	case "warn":
		return logrus.WarnLevel, nil
	case "error":
		return logrus.ErrorLevel, nil
	default:
		return logrus.InfoLevel, nil
	}
}
