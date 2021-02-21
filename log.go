package hlog

import "github.com/sirupsen/logrus"

func Now() (l *logrus.Logger) {
	defer func() {
		l = logrus.StandardLogger()
	}()
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel)
	//logrus.SetReportCaller(true) // 日志显示行号，默认false
	return
}
