package stream

import (
	"net"
	"strconv"
	"time"

	"github.com/dxcweb/go-nat-hole/common"
	"github.com/xtaci/smux"
)

//Stream  流
type Stream struct {
	p    *smux.Stream
	id   uint32
	name string
}

var streams = make([]*Stream, 0, 20)

func GetStreamByName(name string) *Stream {
	for _, value := range streams {
		if value.name == name {
			return value
		}
	}
	return nil
}
func GetStreamByID(s string) *Stream {
	id64, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return nil
	}
	id := uint32(id64)
	for _, value := range streams {
		if value.id == id {
			return value
		}
	}
	return nil
}

//SetName 设置名称
func (s *Stream) SetName(name string) {
	s.name = name
}

//AcceptStream 接受流
func AcceptStream(mux *smux.Session) (*Stream, error) {
	p, err := mux.AcceptStream()
	if err != nil {
		return nil, err
	}
	s := &Stream{
		p:  p,
		id: common.Rand(),
	}
	streams = append(streams, s)
	return s, nil
}

//ID 返回id
func (s *Stream) ID() uint32 {
	return s.id
}

//Close Close
func (s *Stream) Close() error {
	return s.p.Close()
}

// LocalAddr satisfies net.Conn interface
func (s *Stream) LocalAddr() net.Addr {
	return s.p.LocalAddr()
}

func (s *Stream) Read(b []byte) (n int, err error) {
	return s.p.Read(b)
}

// RemoteAddr satisfies net.Conn interface
func (s *Stream) RemoteAddr() net.Addr {
	return s.p.RemoteAddr()
}

// SetDeadline sets both read and write deadlines as defined by
// net.Conn.SetDeadline.
// A zero time value disables the deadlines.
func (s *Stream) SetDeadline(t time.Time) error {
	return s.p.SetDeadline(t)
}

// SetReadDeadline sets the read deadline as defined by
// net.Conn.SetReadDeadline.
// A zero time value disables the deadline.
func (s *Stream) SetReadDeadline(t time.Time) error {
	return s.p.SetReadDeadline(t)
}

// SetWriteDeadline sets the write deadline as defined by
// net.Conn.SetWriteDeadline.
// A zero time value disables the deadline.
func (s *Stream) SetWriteDeadline(t time.Time) error {
	return s.p.SetWriteDeadline(t)
}

// Write implements net.Conn
func (s *Stream) Write(b []byte) (n int, err error) {
	return s.p.Write(b)
}
