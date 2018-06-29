package manage

import (
	"github.com/dxcweb/go-nat-hole/intermediary/stream"
	"github.com/sirupsen/logrus"
)

//Server 服务端
type Server struct {
	name string
	p    *stream.Stream //流
}

var servers = make([]*Server, 0, 10)

//NewServer 新建一个服务
func NewServer(name string, p *stream.Stream) {
	server := &Server{name, p}
	servers = append(servers, server)
	logrus.Info("新建一个服务，当前服务数量", len(servers))
}

func CreateHole() {

}
