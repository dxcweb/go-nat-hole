package intermediaryclient

import (
	"io"

	"github.com/dxcweb/go-nat-hole/common"
)

func sendFindServer(p io.Writer, name string) {
	optBytes := []byte{common.FlagFindServer}
	nameBytes := []byte(name)
	data := make([]byte, len(nameBytes)+1)
	copy(data[0:1], optBytes)
	copy(data[1:], nameBytes)
	p.Write(data)
}
