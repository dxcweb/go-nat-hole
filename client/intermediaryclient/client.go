package intermediaryclient

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/dxcweb/go-nat-hole/client/conf"
	"github.com/dxcweb/go-nat-hole/common"
	"github.com/sirupsen/logrus"
	kcp "github.com/xtaci/kcp-go"
	"github.com/xtaci/smux"
)

//RunIntermediaryClient 运行客户端
func RunIntermediaryClient(config *conf.Config, conn *net.UDPConn, addrChen chan string) error {
	kcpconn, err := common.NewConn(config.Key, config.RemoteAddr, conn)
	// defer kcpconn.Close()
	if err != nil {
		logrus.Error("客户端启动失败：", err)
		return err
	}
	handleClient(kcpconn, config, addrChen)
	return nil
}
func handleClient(conn *kcp.UDPSession, config *conf.Config, addrChen chan string) error {
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
		logrus.Info("关闭和中介服务器的通讯")
		p.Close()
		close(findServerChan)
		time.Sleep(time.Second * 1)
	}()
	go findServer(p, config.ServerName, findServerChan)
	for {
		buf := make([]byte, 1024)
		nr, er := p.Read(buf)
		if nr > 0 {
			reader(p, buf[0:nr], findServerChan, addrChen, sess)
		}
		if er != nil {
			fmt.Println(333)
			break
		}
	}
	return nil
}

func findServer(p io.Writer, name string, exitChan chan int) {
	sendCount := 1
	for {
		select {
		case <-exitChan:
			fmt.Println("aqsdfasdfasd")
			return
		default:
			logrus.Info("发送发现服务,发送次数：", sendCount)
			sendFindServer(p, name)
			sendCount++
			time.Sleep(time.Second * 3)
		}
	}
}
