package main

import (
	"os"

	"github.com/urfave/cli"
)

func main() {
	myApp := cli.NewApp()
	myApp.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "listen,l",
			Value: ":18888",
			Usage: "kcp server listen address",
		},
		cli.StringFlag{
			Name:   "key",
			Value:  "123456",
			Usage:  "客户机与服务器之间的预共享秘密",
			EnvVar: "KCPTUN_KEY",
		},
	}
	myApp.Action = func(c *cli.Context) {
		config := Config{}
		config.Key = c.String("key")
		config.Listen = c.String("listen")
		RunIntermediary(config)
	}
	myApp.Run(os.Args)
}
