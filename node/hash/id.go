package hash

import "crypto/sha1"

type ID struct {
	id []byte
}

func NewID(i int) {
	sha := sha1.New()
	sha.Write([]byte{byte(i)})
}
