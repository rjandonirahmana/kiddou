package base

import (
	"io"
	"math/rand"
	"time"
)

var num = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

func GenerateUserID() string {
	b := make([]byte, 20)
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	n, err := io.ReadAtLeast(random, b, 20)
	if n != 20 {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = num[int(b[i])%len(num)]
	}
	return string(b)
}
