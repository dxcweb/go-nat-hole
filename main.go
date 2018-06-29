package main

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

func main() {
	addr := "255.255.255.255:65535"
	b := bytes.Buffer{}
	b.WriteString(addr)
	b.WriteString("|")
	b.WriteString("asdf")
	// string()
	c := b.Bytes()
	s := string(c)
	fmt.Println(s)
	// fmt.Println([]byte())
	// fmt.Println([]byte("255.255.255.255:1"))
}

func AddrToString(addr []byte) {
	s0 := string(int(addr[0]))
	s1 := string(int(addr[1]))
	s := []string{s0, s1}
	fmt.Println(strings.Join(s, "."))
}

//AddrToBytes 字符串地址转换为字节
func AddrToBytes(addr string) []byte {
	temp := strings.Split(addr, ":")
	parts := strings.Split(temp[0], ".")
	b0, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil
	}
	b1, _ := strconv.Atoi(parts[1])
	b2, _ := strconv.Atoi(parts[2])
	b3, _ := strconv.Atoi(parts[3])

	b := make([]byte, 4)
	b[0] = byte(b0)
	b[1] = byte(b1)
	b[2] = byte(b2)
	b[3] = byte(b3)
	return b
}
