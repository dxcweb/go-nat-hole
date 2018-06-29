package common

import (
	"net"

	"github.com/golang/snappy"
)

//CompStream 流压缩
type CompStream struct {
	conn net.Conn
	w    *snappy.Writer
	r    *snappy.Reader
}

//Read 读
func (c *CompStream) Read(p []byte) (n int, err error) {
	return c.r.Read(p)
}

//Write 写
func (c *CompStream) Write(p []byte) (n int, err error) {
	n, err = c.w.Write(p)
	err = c.w.Flush()
	return n, err
}

//Close 关闭
func (c *CompStream) Close() error {
	return c.conn.Close()
}

// LocalAddr returns the local network address.
func (c *CompStream) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

// LocalAddr returns the local network address.
func (c *CompStream) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

//NewCompStream new一个流压缩
func NewCompStream(conn net.Conn) *CompStream {
	c := new(CompStream)
	c.conn = conn
	c.w = snappy.NewBufferedWriter(conn)
	c.r = snappy.NewReader(conn)
	return c
}
