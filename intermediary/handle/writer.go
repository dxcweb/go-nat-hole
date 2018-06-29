package handle

import (
	"bytes"
	"io"
	"strconv"

	"github.com/dxcweb/go-nat-hole/common"
)

func sendRegisterConfirm(p io.Writer) {
	optBytes := []byte{common.FlagRegisterConfirm}
	p.Write(optBytes)
}

func sendCreateHole(p io.Writer, addr string, id uint32) {
	idStr := strconv.FormatUint(uint64(id), 10)
	optBytes := []byte{common.FlagCreateHole}
	buf := bytes.Buffer{}
	buf.WriteString(addr)
	buf.WriteString("|")
	buf.WriteString(idStr)
	dataBytes := buf.Bytes()
	data := make([]byte, len(dataBytes)+1)
	copy(data[0:1], optBytes)
	copy(data[1:], dataBytes)
	p.Write(data)
}

func sendConnectServer(p io.Writer, addr string) {
	optBytes := []byte{common.FlagConnectServer}
	dataBytes := []byte(addr)
	data := make([]byte, len(dataBytes)+1)
	copy(data[0:1], optBytes)
	copy(data[1:], dataBytes)
	p.Write(data)
}
