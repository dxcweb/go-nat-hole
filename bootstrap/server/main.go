package main

import (
	"os"

	"github.com/dxcweb/go-nat-hole/server"
	"github.com/dxcweb/go-nat-hole/server/conf"
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
			Name:  "remoteaddr, r",
			Value: "127.0.0.1:18888",
			Usage: "intermediary地址",
		},
		cli.StringFlag{
			Name:  "name, n",
			Value: "server",
			Usage: "服务名字，用于客户端指定链接的唯一标识",
		},
	}
	myApp.Action = func(c *cli.Context) {
		config := conf.Config{}
		config.Key = c.String("key")
		config.RemoteAddr = c.String("remoteaddr")
		config.Name = c.String("name")

		server.RunServer(&config)
	}
	myApp.Run(os.Args)
}
