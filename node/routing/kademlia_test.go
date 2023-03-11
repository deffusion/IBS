package routing

import (
	"fmt"
	"testing"
)

func TestLocate(t *testing.T) {
	for i := 1; i < 8; i++ {
		fmt.Printf("peer %d in bucket %d\n", i, locate(uint64(0), uint64(i)))
	}
}

func TestFakeID(t *testing.T) {
	fmt.Println(fakeIDForBucket(0, 63))
}
