package common

import (
	"crypto/sha1"
	"net"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	kcp "github.com/xtaci/kcp-go"
	"github.com/xtaci/smux"
	"golang.org/x/crypto/pbkdf2"
)

var (
	// SALT is use for pbkdf2 key expansion
	SALT        = "go-nat-hole"
	dataShard   = 10
	parityShard = 10
	//SockBuf socket buffer size in bytes
	SockBuf = 4194304
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
	lis, err := kcp.ListenWithOptions(listen, block, dataShard, parityShard)
	if err != nil {
		return nil, err
	}
	// TODO 不知道DSCP是啥
	// if err := lis.SetDSCP(0); err != nil {
	// 	logrus.Warn("SetDSCP:", err)
	// }
	if err := lis.SetReadBuffer(SockBuf); err != nil {
		logrus.Warn("SetReadBuffer:", err)
	}
	if err := lis.SetWriteBuffer(SockBuf); err != nil {
		logrus.Warn("SetWriteBuffer:", err)
	}

	return lis, nil
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
	conn, err := kcp.NewConn(remoteAddr, block, dataShard, parityShard, &connectedUDPConn{udpconn})
	if err != nil {
		return nil, err
	}

	SetConnOption(conn)

	// TODO 不知道DSCP是啥
	// if err := conn.SetDSCP(0); err != nil {
	// 	logrus.Warn("SetDSCP:", err)
	// }
	if err := conn.SetReadBuffer(SockBuf); err != nil {
		logrus.Warn("SetReadBuffer:", err)
	}
	if err := conn.SetWriteBuffer(SockBuf); err != nil {
		logrus.Warn("SetWriteBuffer:", err)
	}

	return conn, nil
}

//SetConnOption 设置连接选项
func SetConnOption(conn *kcp.UDPSession) {
	conn.SetStreamMode(true)
	conn.SetWriteDelay(true)
	NoDelay, Interval, Resend, NoCongestion := 1, 10, 2, 1
	conn.SetNoDelay(NoDelay, Interval, Resend, NoCongestion)
	conn.SetMtu(1350)              // 设置UDP数据包的最大传输单元
	conn.SetWindowSize(1024, 1024) //  设置发送的窗口大小
	conn.SetACKNoDelay(false)      //当接收到数据包时立即刷新ACK
}

//SmuxConfig 多路复用配置
func SmuxConfig() *smux.Config {
	smuxConfig := smux.DefaultConfig()
	smuxConfig.MaxReceiveBuffer = SockBuf
	smuxConfig.KeepAliveInterval = time.Duration(10) * time.Second
	return smuxConfig
}

type connectedUDPConn struct{ *net.UDPConn }

// WriteTo redirects all writes to the Write syscall, which is 4 times faster.
func (c *connectedUDPConn) WriteTo(b []byte, addr net.Addr) (int, error) {
	return c.Write(b)
}
