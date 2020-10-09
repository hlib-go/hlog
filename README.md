# hlog

#### 介绍
logger


使用说明：
```golang

var log = hlof.New()

// 日志格式为JSON
logrus.SetFormatter(&logrus.JSONFormatter{})

// 设置日志级别
logrus.SetLevel(logrus.InfoLevel)

// 使用阿里云SLS日志Hook
hlog.AddAliyunSlsHook(sls)

```