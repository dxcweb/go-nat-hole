package common

import (
	"math/rand"
	"time"
)

var r = rand.New(rand.NewSource(time.Now().Unix()))

//Rand 生产uint32随机数
func Rand() uint32 {
	return r.Uint32()
}
