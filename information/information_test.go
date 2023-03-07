package information

import (
	"fmt"
	"sort"
	"testing"
)

func TestSort(t *testing.T) {
	p1 := NewPacket(1, 1024, nil, nil, nil, 10, nil)
	p2 := NewPacket(1, 1024, nil, nil, nil, 20, nil)
	p3 := NewPacket(1, 1024, nil, nil, nil, 5, nil)
	p4 := NewPacket(1, 1024, nil, nil, nil, 100, nil)
	p5 := NewPacket(1, 1024, nil, nil, nil, 30, nil)
	p6 := NewPacket(1, 1024, nil, nil, nil, 11, nil)
	var ps = Packets{p1, p2, p3, p4, p5, p6}
	sort.Sort(ps)
	for _, p := range ps {
		fmt.Println(p)
	}
}
