package handle

import (
	"github.com/dxcweb/go-nat-hole/common"
	"github.com/dxcweb/go-nat-hole/intermediary/stream"
	"github.com/sirupsen/logrus"
)

//Reader 解析接受到的消息
func Reader(p *stream.Stream, buf []byte) {
	switch buf[0] {
	case common.FlagRegister:
		name := string(buf[1:])
		p.SetName(name)
		sendRegisterConfirm(p)
		logrus.Info("接受到注册事件,来自", name, ",nat：", p.RemoteAddr())
		break
	case common.FlagFindServer:
		name := string(buf[1:])
		s := stream.GetStreamByName(name)
		if s == nil {
			logrus.Warn("接受到客户端：", p.RemoteAddr(), ",发现服务请求失败!,服务名：", name)
		} else {
			sendCreateHole(s, p.RemoteAddr().String(), p.ID())
			logrus.Info("接受到客户端：", p.RemoteAddr(), ",发现服务请求：", name)
		}
		break
	case common.FlagCreateHoleFinish:
		clientID := string(buf[1:])
		s := stream.GetStreamByID(clientID)
		if s == nil {
			logrus.Warn("创建打洞服务：", p.RemoteAddr(), "找不到客户端：", clientID)
		} else {
			sendConnectServer(s, p.RemoteAddr().String())
			logrus.Info("创建打洞服务：", p.RemoteAddr())
		}
		break
	default:
		p.Close()
		logrus.Error("接受到错误操作", buf)
	}
}
