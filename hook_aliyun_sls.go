package logger

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/aliyun-log-go-sdk/producer"
	"github.com/sirupsen/logrus"
	"net"
	"time"
)

// Hook 阿里云SLS存储日志
func AddAliyunSlsHook(sls *AliyunSls) {
	logrus.AddHook(&AliyunSlsHook{AliyunSls: sls})
	logrus.Info("Add AliyunSlsHook success...")
}

type AliyunSlsHook struct {
	AliyunSls *AliyunSls
}

func (hook *AliyunSlsHook) Fire(entry *logrus.Entry) error {
	bytes, err := entry.Bytes()
	if err != nil {
		return err
	}
	// 转换为Map，存入SLS
	var mkv map[string]string
	err = json.Unmarshal(bytes, &mkv)
	if err != nil {
		return err
	}
	hook.AliyunSls.SendLog(mkv)
	return nil
}

func (hook *AliyunSlsHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

type AliyunSls struct {
	// 通过阿里云控制台获取配置信息
	Endpoint        string `json:"endpoint"`
	AccessKeyID     string `json:"accessKeyId"`
	AccessKeySecret string `json:"accessKeySecret"`
	Project         string `json:"project"`
	Logstore        string `json:"logstore"`
	Topic           string `json:"topic"`

	producer *producer.Producer `json:"-"`
	source   string             `json:"-"`
}

func (sls *AliyunSls) Producer() *producer.Producer {
	if sls.producer != nil {
		return sls.producer
	}
	sls.source = getLocalIP()

	producerConfig := producer.GetDefaultProducerConfig()
	producerConfig.Endpoint = sls.Endpoint
	producerConfig.AccessKeyID = sls.AccessKeyID
	producerConfig.AccessKeySecret = sls.AccessKeySecret
	producer := producer.InitProducer(producerConfig)
	producer.Start() // 启动producer实例
	sls.producer = producer
	return sls.producer
}

func (sls *AliyunSls) SendLog(kv map[string]string) {
	log := producer.GenerateLog(uint32(time.Now().Unix()), kv)
	// himkt   msd   127.0.0.1   k-v
	err := sls.Producer().SendLog(sls.Project, sls.Logstore, sls.Topic, sls.source, log)
	if err != nil {
		fmt.Println("SendLog Error:", err.Error())
	}
}

// 获取本机网卡IP
func getLocalIP() (ipv4 string) {
	var (
		err     error
		addrs   []net.Addr
		addr    net.Addr
		ipNet   *net.IPNet // IP地址
		isIpNet bool
	)
	// 获取所有网卡
	if addrs, err = net.InterfaceAddrs(); err != nil {
		return
	}
	// 取第一个非lo的网卡IP
	for _, addr = range addrs {
		// 这个网络地址是IP地址: ipv4, ipv6
		if ipNet, isIpNet = addr.(*net.IPNet); isIpNet && !ipNet.IP.IsLoopback() {
			// 跳过IPV6
			if ipNet.IP.To4() != nil {
				ipv4 = ipNet.IP.String() // 192.168.1.1
				return
			}
		}
	}

	if err != nil {
		ipv4 = "127.0.0.1"
	}
	return
}
