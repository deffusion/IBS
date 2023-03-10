package num_set

import (
	"fmt"
	"testing"
)

func TestSet(t *testing.T) {
	s := NewSet()
	//s.Insert(uint64(3))
	s.Insert(uint64(5))
	s.Insert(uint64(2))
	s.Insert(uint64(8))
	//s.Insert(uint64(9))
	//s.Insert(uint64(0))
	//s.Insert(uint64(4))
	//s.Insert(uint64(4))
	//s.Insert(uint64(4))
	//s.Insert(uint64(4))
	//s.Insert(uint64(1))
	//s.Insert(uint64(7))
	//s.Insert(uint64(6))
	//fmt.Println(s.s)
	//fmt.Println("locate:", s.locate(3))
	fmt.Println(s.Around(3, 10))
}
