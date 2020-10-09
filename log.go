package logger

import "github.com/sirupsen/logrus"

func Now() (l *logrus.Logger) {
	defer func() {
		l = logrus.StandardLogger()
	}()
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel)
	return
}
