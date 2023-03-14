package node

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"testing"
)

func TestSha1(t *testing.T) {
	sha := sha1.New()
	sha.Write([]byte{'1'})
	hash := sha.Sum(nil)

	fmt.Println(hex.EncodeToString(hash))
}
func TestMap(t *testing.T) {
	m := make(map[int]int)
	m[0] = 1
	b, ok := m[0]
	fmt.Println(b, ok)
}
func TestNodeType(t *testing.T) {
	//var n1 Node = NewNeNode(1, 1, 1, "cn", nil)
	var n2 Node = NewBasicNode(2, 2, 2, "cn", nil)
	switch n2.(type) {
	case *NeNode:
		fmt.Println("ne")
	case *BasicNode:
		fmt.Println("basic")
	default:
		fmt.Println("none")
	}
}
