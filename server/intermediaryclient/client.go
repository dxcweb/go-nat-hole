package intermediaryclient

import (
	"io"
	"time"

	"github.com/dxcweb/go-nat-hole/common"
	"github.com/dxcweb/go-nat-hole/server/conf"
	"github.com/sirupsen/logrus"
	kcp "github.com/xtaci/kcp-go"
	"github.com/xtaci/smux"
)

//RunIntermediaryClient 运行中介服务器客户端
func RunIntermediaryClient(config *conf.Config) error {
	conn, err := common.UDPClientSimple(config.Key, config.RemoteAddr)
	if err != nil {
		logrus.Error("服务端启动失败：", err)
		return err
	}
	handleClient(conn, config)
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
	registerChan := make(chan int, 1)
	defer func() {
		p.Close()
		close(registerChan)
		time.Sleep(time.Second * 1)
	}()
	go register(p, config.Name, registerChan)
	for {
		buf := make([]byte, 1024)
		nr, er := p.Read(buf)
		if nr > 0 {
			reader(p, buf[0:nr], registerChan, config)
		}
		if er != nil {
			// if er != io.EOF {
			// 	logrus.Info("错误", er)
			// }
			break
		}
	}
	return nil
}

func register(p io.Writer, name string, exitChan chan int) {
	sendCount := 1
	for {
		select {
		case <-exitChan:
			return
		default:
			logrus.Info("发送注册消息,发送次数：", sendCount)
			sendRegister(p, name)
			sendCount++
			time.Sleep(time.Second * 3)
		}
	}
}
