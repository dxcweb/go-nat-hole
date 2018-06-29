package hole

import (
	"fmt"
	"net"
	"time"

	"github.com/dxcweb/go-nat-hole/common"
	"github.com/dxcweb/go-nat-hole/intermediary/stream"
	"github.com/dxcweb/go-nat-hole/server/conf"
	"github.com/sirupsen/logrus"
	"github.com/xtaci/smux"
)

func CreateHole(config *conf.Config, clientAddr, clientID string) {
	conn := newUdpconn(config.Key)
	pc, _ := openStream(config.Key, clientAddr, conn)
	pi, _ := openStream(config.Key, config.RemoteAddr, conn)
	pc.Write([]byte{1})
	pc.Close()
	sendCreateHoleFinish(pi, clientID)
	pi.Close()
	time.Sleep(time.Second * 2)
}

func openStream(key, raddr string, conn *net.UDPConn) (*smux.Stream, error) {
	kcpconn, err := common.NewConn(key, raddr, conn)
	if err != nil {
		logrus.Error("服务端启动失败：", err)
		return nil, err
	}
	smuxConfig := common.SmuxConfig()
	sess, err := smux.Client(common.NewCompStream(kcpconn), smuxConfig)
	if err != nil {
		logrus.Error("smux.Client失败：", err)
		return nil, err
	}
	p, err := sess.OpenStream()
	if err != nil {
		logrus.Error("OpenStream失败：", err)
		return nil, err
	}
	return p, nil
}
func newServer(key string, conn *net.UDPConn) error {
	lis, err := common.UDPServerByConn(key, conn)
	if err != nil {
		logrus.Error("监听UDP端口失败：", err)
		return err
	}
	logrus.Info("中aaaa启动成功端口为:")
	for {
		if conn, err := lis.AcceptKCP(); err == nil {
			common.SetConnOption(conn)
			go handleMux(common.NewCompStream(conn))
		} else {
			logrus.Error("lis.AcceptKCP错误", err)
		}
	}
}

// 多路复用
func handleMux(conn *common.CompStream) {
	fmt.Println(123123)
	smuxConfig := common.SmuxConfig()
	mux, err := smux.Server(conn, smuxConfig)
	if err != nil {
		logrus.Error("启动多路复用服务失败", err)
		return
	}
	defer mux.Close()
	for {
		p, err := stream.AcceptStream(mux)
		if err != nil {
			logrus.Error("错误123")
			return
		}
		fmt.Println(31)
		go func() {
			defer p.Close()
			buf := make([]byte, 32*1024)
			for {
				n, er := p.Read(buf)
				if n > 0 {
					fmt.Println("接受到消息:", buf[0:n])
				}
				if er != nil {
					// if er != io.EOF {
					// 	logrus.Info("错误", er)
					// }
					break
				}
			}
		}()
	}
}

var temp *net.UDPConn

func newUdpconn(key string) *net.UDPConn {
	if temp == nil {
		localAddr := ":19102"
		conn, _ := common.UDPconn(localAddr)
		go newServer(key, conn)
		temp = conn
	}
	return temp
}
