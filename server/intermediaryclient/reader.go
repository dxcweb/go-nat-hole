package intermediaryclient

import (
	"io"
	"strings"

	"github.com/dxcweb/go-nat-hole/common"
	"github.com/dxcweb/go-nat-hole/server/conf"
	"github.com/dxcweb/go-nat-hole/server/hole"
	"github.com/sirupsen/logrus"
)

func reader(p io.WriteCloser, buf []byte, registerChan chan int, config *conf.Config) {
	switch buf[0] {
	case common.FlagRegisterConfirm:
		logrus.Info("注册成功，等待接受命令")
		registerChan <- 1
		break
	case common.FlagCreateHole:
		data := string(buf[1:])
		dataArr := strings.Split(data, "|")
		hole.CreateHole(config, dataArr[0], dataArr[1])
		break
	default:
		p.Close()
		logrus.Error("接受到错误操作", buf)
	}
}
