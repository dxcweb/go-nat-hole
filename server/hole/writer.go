package hole

import (
	"io"

	"github.com/dxcweb/go-nat-hole/common"
)

func sendCreateHoleFinish(p io.Writer, clientID string) {
	optBytes := []byte{common.FlagCreateHoleFinish}
	dataBytes := []byte(clientID)
	data := make([]byte, len(dataBytes)+1)
	copy(data[0:1], optBytes)
	copy(data[1:], dataBytes)
	p.Write(data)
}
