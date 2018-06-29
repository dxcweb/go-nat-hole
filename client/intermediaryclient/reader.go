package intermediaryclient

import (
	"github.com/dxcweb/go-nat-hole/common"
	"github.com/sirupsen/logrus"
	"github.com/xtaci/smux"
)

//Reader 解析接受到的消息
func reader(p *smux.Stream, buf []byte, findServerChan chan int, addrChen chan string, s *smux.Session) {
	switch buf[0] {
	case common.FlagConnectServer:
		addr := string(buf[1:])
		findServerChan <- 1
		addrChen <- addr
		p.Close()
		break
	default:
		p.Close()
		logrus.Error("接受到错误操作", buf)
	}
}
