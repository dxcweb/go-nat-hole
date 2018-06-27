package main

import (
	"fmt"
	"io"

	"github.com/dxcweb/go-nat-hole/common"
	"github.com/sirupsen/logrus"
	"github.com/xtaci/smux"
)

//RunIntermediary 运行中介服务器
func RunIntermediary(config Config) error {
	lis, err := common.UDPServer(config.Key, config.Listen)
	if err != nil {
		logrus.Error("监听UDP端口失败：", err)
		return err
	}
	logrus.Info("中介服务器启动成功端口为:", config.Listen)
	for {
		if conn, err := lis.AcceptKCP(); err == nil {
			logrus.Info("客户端链接地址为:", conn.RemoteAddr())
			common.SetConnOption(conn)
			go handleMux(common.NewCompStream(conn), &config)
		} else {
			logrus.Error("lis.AcceptKCP错误", err)
		}
	}
}

// 多路复用
func handleMux(conn io.ReadWriteCloser, config *Config) {
	smuxConfig := common.SmuxConfig()

	mux, err := smux.Server(conn, smuxConfig)
	if err != nil {
		logrus.Error("启动多路复用服务失败", err)
		return
	}
	defer mux.Close()
	for {
		fmt.Println(123)
		p, err := mux.AcceptStream()
		if err != nil {
			logrus.Error("接受流:", err)
			return
		}
		go func() {
			buf := make([]byte, 32*1024)
			for {
				nr, er := p.Read(buf)
				if nr > 0 {
					fmt.Println("接受到数据", nr)
				}
				if er != nil {
					if er != io.EOF {
						logrus.Info("EOF", err)
					}
					break
				}
			}
			logrus.Info("退出读", err)
		}()
	}
}
