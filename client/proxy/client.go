package proxy

import (
	"fmt"
	"net"
	"time"

	"github.com/dxcweb/go-nat-hole/client/conf"
	"github.com/dxcweb/go-nat-hole/common"
	"github.com/sirupsen/logrus"
	kcp "github.com/xtaci/kcp-go"
	"github.com/xtaci/smux"
)

func NewProxyClient(config *conf.Config, remoteAddr string, conn *net.UDPConn) error {
	kcpconn, err := common.NewConn(config.Key, remoteAddr, conn)
	defer conn.Close()
	if err != nil {
		logrus.Error("客户端启动失败：", err)
		return err
	}
	handleClient(kcpconn, config)
	return nil
}

func handleClient(conn *kcp.UDPSession, config *conf.Config) error {
	smuxConfig := common.SmuxConfig()
	sess, err := smux.Client(common.NewCompStream(conn), smuxConfig)
	if err != nil {
		logrus.Error("smux.Client失败：", err)
		return err
	}
	p, err := sess.OpenStream()
	if err != nil {
		logrus.Error("OpenStream失败：", err)
		return err
	}
	findServerChan := make(chan int)
	defer func() {
		logrus.Info("关闭和服务端的通讯")
		p.Close()
		close(findServerChan)
		time.Sleep(time.Second * 1)
	}()
	p.Write([]byte{3, 3, 3})
	for {
		buf := make([]byte, 1024)
		nr, er := p.Read(buf)
		if nr > 0 {
			fmt.Println("123", buf[:nr])
		}
		if er != nil {
			fmt.Println(333)
			break
		}
	}
	return nil
}
