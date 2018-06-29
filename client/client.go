package client

import (
	"github.com/dxcweb/go-nat-hole/client/conf"
	"github.com/dxcweb/go-nat-hole/client/intermediaryclient"
	"github.com/dxcweb/go-nat-hole/client/proxy"
	"github.com/dxcweb/go-nat-hole/common"
	"github.com/sirupsen/logrus"
)

//RunClient 运行客户端
func RunClient(config *conf.Config) {
	conn, _ := common.UDPconn(config.LocalAddr)
	// 获取连接中介服务器获取目标地址
	addr := intermediaryclient.GetServerAddr(config, conn)
	logrus.Info("获得到服务端地址", addr)
	proxy.NewProxyClient(config, addr, conn)
}
