package main

import (
	"os"

	"github.com/dxcweb/go-nat-hole/client"
	"github.com/dxcweb/go-nat-hole/client/conf"
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
			Value: "127.0.0.1:18888",
			Usage: "intermediary地址",
		},
		cli.StringFlag{
			Name:  "servername, s",
			Value: "server",
			Usage: "服务端名称",
		},
	}
	myApp.Action = func(c *cli.Context) {
		config := conf.Config{}
		config.Key = c.String("key")
		config.LocalAddr = c.String("localaddr")
		config.RemoteAddr = c.String("remoteaddr")
		config.ServerName = c.String("servername")

		client.RunClient(&config)
	}
	myApp.Run(os.Args)
}
