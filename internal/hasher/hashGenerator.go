package hasher

import (
	"crypto/sha1"
	"fmt"
)

const salt = "gqgwgd1g21gehwdwh08w7dbb1y2hshsdasd"

func GeneratePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
