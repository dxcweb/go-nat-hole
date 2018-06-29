package server

import (
	"time"

	"github.com/dxcweb/go-nat-hole/server/conf"
	"github.com/dxcweb/go-nat-hole/server/intermediaryclient"
	"github.com/sirupsen/logrus"
)

//RunServer 运行服务端
func RunServer(config *conf.Config) {
	for {
		intermediaryclient.RunIntermediaryClient(config)
		logrus.Warn("连接到中介服务器失败，正在重试...\n\n\n")
		time.Sleep(time.Second * 3)
	}

	// hole.CreateHole(config)
	// time.Sleep(time.Second * 3)
}
