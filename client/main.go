package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dxcweb/go-nat-hole/common"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func main() {
	myApp := cli.NewApp()
	myApp.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "key",
			Value:  "123456",
			Usage:  "客户机与服务器之间的预共享秘密",
			EnvVar: "KCPTUN_KEY",
		},
		cli.StringFlag{
			Name:  "localaddr,l",
			Value: ":19900",
			Usage: "UDP监听的本机地址",
		},
		cli.StringFlag{
			Name:  "remoteaddr, r",
			Value: "127.0.0.1:29900",
			Usage: "intermediary地址",
		},
	}
	myApp.Action = func(c *cli.Context) error {
		config := Config{}
		config.Key = c.String("key")
		config.LocalAddr = c.String("localaddr")
		config.RemoteAddr = c.String("remoteaddr")

		conn, err := common.UDPClient(config.Key, config.LocalAddr, config.RemoteAddr)
		if err != nil {
			logrus.Error("UDP客户端启动失败：", err)
			return err
		}
		_, err = conn.Write([]byte{1, 2})
		fmt.Println("err", err)
		for {
			data := make([]byte, 1024)
			n, err := conn.Read(data)
			if err != nil {
				log.Printf("error during read: %s\n", err)
			} else {
				log.Printf("收到数据:%s\n", data[:n])
			}
		}
	}
	myApp.Run(os.Args)
}
