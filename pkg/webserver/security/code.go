package security

import (
	"math/rand"
	"time"
)

func Code() string {
	var CodeSet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	code := ""

	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 6; i++ {
		code += string(CodeSet[rand.Intn(len(CodeSet))])
	}

	return code
}
