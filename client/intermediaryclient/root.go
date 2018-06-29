package intermediaryclient

import (
	"net"
	"time"

	"github.com/dxcweb/go-nat-hole/client/conf"
	"github.com/sirupsen/logrus"
)

//GetServerAddr 获取服务端地址
func GetServerAddr(config *conf.Config, conn *net.UDPConn) string {
	addrChen := make(chan string, 1)
	i := 0
	for {
		select {
		case addr := <-addrChen:
			return addr
		default:
			if i > 0 {
				logrus.Warn("连接到中介服务器失败，正在重试...\n\n\n")
				time.Sleep(time.Second * 3)
			}
			RunIntermediaryClient(config, conn, addrChen)
			i++
		}
	}
}
