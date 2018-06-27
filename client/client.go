package main

import (
	"time"

	"github.com/dxcweb/go-nat-hole/common"
	"github.com/sirupsen/logrus"
	kcp "github.com/xtaci/kcp-go"
	"github.com/xtaci/smux"
)

//RunClient 运行客户端
func RunClient(config Config) error {
	conn, err := common.UDPClient(config.Key, config.LocalAddr, config.RemoteAddr)
	if err != nil {
		logrus.Error("客户端启动失败：", err)
		return err
	}
	logrus.Info("客户端启动成功，本地：", config.LocalAddr, "远程：", config.RemoteAddr)
	handleClient(conn)
	return nil
}
func handleClient(conn *kcp.UDPSession) error {
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
	defer func() {
		err := p.Close()
		time.Sleep(time.Second * 3)
		logrus.Info("关闭", err)
	}()
	b := make([]byte, 64*1024)
	p.Write(b)
	//
	// go func() {
	// 	buf := make([]byte, 32*1024)
	// 	for {
	// 		nr, er := p.Read(buf)
	// 		if nr > 0 {
	// 			fmt.Println("接受到数据", buf[0:nr])
	// 		}
	// 		if er != nil {
	// 			if er != io.EOF {
	// 				logrus.Info("EOF", err)
	// 			}
	// 			break
	// 		}
	// 	}
	// }()
	return nil
}
