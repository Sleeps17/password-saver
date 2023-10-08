package random

import (
	"math/rand"
	"time"
)

func GeneratePassword(lenght int) string {

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "abcdefghigklmnopqrstuvwxyz" + "0123456789" + "/><.-*&^$#!@~")

	result := make([]rune, lenght)

	for i := 0; i < lenght; i++ {
		result[i] = chars[rnd.Intn(len(chars))]
	}

	return string(result)
}
