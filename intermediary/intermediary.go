package intermediary

import (
	"fmt"

	"github.com/dxcweb/go-nat-hole/common"
	"github.com/dxcweb/go-nat-hole/intermediary/conf"
	"github.com/dxcweb/go-nat-hole/intermediary/handle"
	"github.com/dxcweb/go-nat-hole/intermediary/stream"
	"github.com/sirupsen/logrus"
	"github.com/xtaci/smux"
)

//RunIntermediary 运行中介服务器
func RunIntermediary(config conf.Config) error {
	lis, err := common.UDPServer(config.Key, config.Listen)
	if err != nil {
		logrus.Error("监听UDP端口失败：", err)
		return err
	}
	logrus.Info("中介服务器启动成功端口为:", config.Listen)
	for {
		if conn, err := lis.AcceptKCP(); err == nil {
			common.SetConnOption(conn)
			go handleMux(common.NewCompStream(conn), &config)
		} else {
			logrus.Error("lis.AcceptKCP错误", err)
		}
	}
}

// 多路复用
func handleMux(conn *common.CompStream, config *conf.Config) {
	fmt.Println(1234)
	smuxConfig := common.SmuxConfig()

	mux, err := smux.Server(conn, smuxConfig)
	if err != nil {
		logrus.Error("启动多路复用服务失败", err)
		return
	}
	defer mux.Close()
	for {
		p, err := stream.AcceptStream(mux)
		fmt.Println("aaaaa")
		if err != nil {
			logrus.Error("错误123")
			return
		}
		go func() {
			defer p.Close()
			buf := make([]byte, 32*1024)
			for {
				n, er := p.Read(buf)
				if n > 0 {
					handle.Reader(p, buf[0:n])
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
