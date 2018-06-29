package client

import (
	"github.com/dxcweb/go-nat-hole/client/conf"
	"github.com/dxcweb/go-nat-hole/client/intermediaryclient"
	"github.com/dxcweb/go-nat-hole/client/proxy"
	"github.com/dxcweb/go-nat-hole/common"
)

//RunClient 运行客户端
func RunClient(config *conf.Config) {
	conn, _ := common.UDPconn(config.LocalAddr)
	// 获取连接中介服务器获取目标地址
	addr := intermediaryclient.GetServerAddr(config, conn)
	proxy.NewProxyClient(config, addr, conn)
}
