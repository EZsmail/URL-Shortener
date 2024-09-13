package random

import (
	"fmt"
	"math/rand"
	"time"
)

var alphabet = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}

func NewRandomString(length int) string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	var alias string

	var test = []rune("iamdiadijaisdjajdiajdsijaisd")

	fmt.Println(test)

	for i := 0; i < length; i++ {
		alias += alphabet[rnd.Intn(26)]
	}

	return alias
}
