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

func UDPconn(localAddr string) (*net.UDPConn, error) {
	laddr, err := net.ResolveUDPAddr("udp", localAddr)
	if err != nil {
		return nil, errors.Wrap(err, "net.ResolveUDPAddr")
	}
	return net.ListenUDP("udp", laddr)
}

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

//UDPServerByConn 启动一个UDP服务
func UDPServerByConn(key string, conn *net.UDPConn) (*kcp.Listener, error) {
	block := GetBlockCrypt(key)
	lis, err := kcp.ServeConn(block, dataShard, parityShard, &connectedUDPConn{conn})
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
	var laddr *net.UDPAddr
	if localAddr != "<nil>" {
		laddr, err = net.ResolveUDPAddr("udp", localAddr)
		if err != nil {
			return nil, errors.Wrap(err, "net.ResolveUDPAddr")
		}
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

func NewConn(key, raddr string, conn *net.UDPConn) (*kcp.UDPSession, error) {
	block := GetBlockCrypt(key)
	// raddr := "127.0.0.1:18888"
	kcpconn, err := kcp.NewConn(raddr, block, dataShard, parityShard, conn)

	if err != nil {
		return nil, err
	}

	SetConnOption(kcpconn)

	// TODO 不知道DSCP是啥
	// if err := conn.SetDSCP(0); err != nil {
	// 	logrus.Warn("SetDSCP:", err)
	// }
	if err := kcpconn.SetReadBuffer(SockBuf); err != nil {
		logrus.Warn("SetReadBuffer:", err)
	}
	if err := kcpconn.SetWriteBuffer(SockBuf); err != nil {
		logrus.Warn("SetWriteBuffer:", err)
	}

	return kcpconn, nil
}

//UDPClientSimple 启动一个UDP客户端
func UDPClientSimple(key, remoteAddr string) (*kcp.UDPSession, error) {
	return UDPClient(key, "<nil>", remoteAddr)
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
