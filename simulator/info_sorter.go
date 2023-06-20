package main

import (
	"IBS/information"
	"fmt"
)

//var mutex sync.Mutex

type PacketSorter struct {
	minHeap []information.Packet
}

//type PacketSorter struct {
//	minHeap []information.Packet
//	length  int
//}
//
//func (s PacketSorter) Length() int {
//	return s.length - 1
//}
//
//func NewInfoSorter() *PacketSorter {
//	s := &PacketSorter{
//		make([]information.Packet, 1, 1<<16),
//		1,
//	}
//	//s.minHeap = append(s.minHeap, nil)
//	s.minHeap[0] = nil
//	return s
//}
//

func NewInfoSorter() *PacketSorter {
	s := &PacketSorter{
		make([]information.Packet, 0, 1<<16),
	}
	return s
}

func (s *PacketSorter) Len() int {
	return len(s.minHeap)
}
func (s *PacketSorter) Less(i, j int) bool {
	return s.minHeap[i].Timestamp() < s.minHeap[j].Timestamp()
}
func (s *PacketSorter) Swap(i, j int) {
	s.minHeap[i], s.minHeap[j] = s.minHeap[j], s.minHeap[i]
}
func (s *PacketSorter) Push(x interface{}) {
	s.minHeap = append(s.minHeap, x.(information.Packet))
}
func (s *PacketSorter) Pop() (v interface{}) {
	v = s.minHeap[len(s.minHeap)-1]
	s.minHeap = s.minHeap[:len(s.minHeap)-1]
	return
}

func (s *PacketSorter) print() {
	for i := 0; i < s.Len(); i++ {
		fmt.Print(s.minHeap[i].Timestamp(), " ")
	}
	fmt.Println()
}

//
//func (s *PacketSorter) Append(info information.Packet) {
//	//mutex.Lock()
//	s.minHeap = append(s.minHeap, info)
//	s.length++
//	//if s.length > MAX {
//	//	MAX = s.length
//	//	fmt.Println("max:", MAX)
//	//}
//	//mutex.Unlock()
//	s.adjustUp(s.length - 1)
//}
//
//func (s *PacketSorter) Take() (information.Packet, error) {
//	//mutex.Lock()
//	if s.length <= 1 {
//		return nil, errors.New("take info from an empty sorter")
//	}
//	toTake := s.minHeap[1]
//	s.minHeap[1], s.minHeap[s.length-1] = s.minHeap[s.length-1], s.minHeap[1]
//	s.length--
//	s.minHeap = s.minHeap[:s.length]
//	//mutex.Unlock()
//	s.adjustDown()
//	return toTake, nil
//}
//
//func (s *PacketSorter) adjustDown() {
//	i := 1
//	for i <= (s.length-1)/2 {
//		// left child default
//		min := 2 * i
//		// if right child exists and bigger
//		if s.length > min+1 && s.minHeap[min+1].Timestamp() < s.minHeap[min].Timestamp() {
//			min += 1
//		}
//		if s.minHeap[i].Timestamp() > s.minHeap[min].Timestamp() {
//			s.minHeap[min], s.minHeap[i] = s.minHeap[i], s.minHeap[min] // swap
//		}
//		i = min
//	}
//}
//
//func (s *PacketSorter) adjustUp(i int) {
//	for ; i >= 2; i = i / 2 {
//		if s.minHeap[i].Timestamp() < s.minHeap[i/2].Timestamp() {
//			s.minHeap[i], s.minHeap[i/2] = s.minHeap[i/2], s.minHeap[i]
//		} else {
//			return
//		}
//	}
//}
