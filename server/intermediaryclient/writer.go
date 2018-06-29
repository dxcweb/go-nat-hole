package intermediaryclient

import (
	"io"

	"github.com/dxcweb/go-nat-hole/common"
)

func sendRegister(p io.Writer, name string) {
	optBytes := []byte{common.FlagRegister}
	nameBytes := []byte(name)
	data := make([]byte, len(nameBytes)+1)
	copy(data[0:1], optBytes)
	copy(data[1:], nameBytes)
	p.Write(data)
}

func SendCreateHoleFinish(p io.Writer, clientID string) {
	optBytes := []byte{common.FlagCreateHoleFinish}
	dataBytes := []byte(clientID)
	data := make([]byte, len(dataBytes)+1)
	copy(data[0:1], optBytes)
	copy(data[1:], dataBytes)
	p.Write(data)
}
