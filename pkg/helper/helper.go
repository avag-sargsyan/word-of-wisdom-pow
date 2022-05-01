package helper

import (
	"fmt"
	"math/rand"
	"time"
)

// For generating random string by the given length
func RandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[:length]
}
