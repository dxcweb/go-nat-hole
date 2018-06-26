package main

import (
	"os"

	"github.com/dxcweb/go-nat-hole/common"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func main() {
	myApp := cli.NewApp()
	myApp.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "listen,l",
			Value: ":29900",
			Usage: "kcp server listen address",
		},
		cli.StringFlag{
			Name:   "key",
			Value:  "123456",
			Usage:  "客户机与服务器之间的预共享秘密",
			EnvVar: "KCPTUN_KEY",
		},
	}
	myApp.Action = func(c *cli.Context) error {
		config := Config{}
		config.Key = c.String("key")
		config.Listen = c.String("listen")

		lis, err := common.UDPServer(config.Key, config.Listen)
		if err != nil {
			logrus.Error("监听UDP端口失败：", err)
			return err
		}
		logrus.Info("监听UDP端口:", config.Listen)
		if conn, err := lis.AcceptKCP(); err == nil {
			logrus.Info("客户端链接地址为:", conn.RemoteAddr())
			for {
				buf := make([]byte, 1024)
				n, err := conn.Read(buf)
				if err != nil {
					logrus.Error("读取内容失败：", err)
				} else {
					logrus.Info("收到数据", buf[:n])
				}
			}
		} else {
			return err
		}
	}
	myApp.Run(os.Args)
}
