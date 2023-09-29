package app

import (
	"github.com/sirupsen/logrus"
	"os"
)

func SetLogrus(level string) {
	logrusLevel, err := logrus.ParseLevel(level)
	if err != nil {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrusLevel)
	}

	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2001-12-31 15:04:05",
	})

	logrus.SetOutput(os.Stdout)
}
