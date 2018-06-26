package common

import (
	"crypto/sha1"
	"net"

	"github.com/pkg/errors"
	kcp "github.com/xtaci/kcp-go"
	"golang.org/x/crypto/pbkdf2"
)

var (
	// SALT is use for pbkdf2 key expansion
	SALT        = "go-nat-hole"
	dataShard   = 10
	parityShard = 10
)

//GetBlockCrypt NewAESBlockCrypt
func GetBlockCrypt(key string) kcp.BlockCrypt {
	pass := pbkdf2.Key([]byte(key), []byte(SALT), 4096, 32, sha1.New)
	block, _ := kcp.NewAESBlockCrypt(pass[:16])
	return block
}

//UDPServer 启动一个UDP服务
func UDPServer(key, listen string) (*kcp.Listener, error) {
	block := GetBlockCrypt(key)
	return kcp.ListenWithOptions(listen, block, dataShard, parityShard)
}

//UDPClient 启动一个UDP客户端
func UDPClient(key, localAddr, remoteAddr string) (*kcp.UDPSession, error) {
	block := GetBlockCrypt(key)

	raddr, err := net.ResolveUDPAddr("udp", remoteAddr)
	if err != nil {
		return nil, errors.Wrap(err, "net.ResolveUDPAddr")
	}
	laddr, err := net.ResolveUDPAddr("udp", localAddr)
	if err != nil {
		return nil, errors.Wrap(err, "net.ResolveUDPAddr")
	}
	udpconn, err := net.DialUDP("udp", laddr, raddr)
	if err != nil {
		return nil, errors.Wrap(err, "net.DialUDP")
	}
	return kcp.NewConn(remoteAddr, block, dataShard, parityShard, &connectedUDPConn{udpconn})
}

type connectedUDPConn struct{ *net.UDPConn }

// WriteTo redirects all writes to the Write syscall, which is 4 times faster.
func (c *connectedUDPConn) WriteTo(b []byte, addr net.Addr) (int, error) {
	return c.Write(b)
}
