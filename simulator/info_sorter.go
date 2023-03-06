package main

import (
	"IBS/information"
	"errors"
	"fmt"
)

//var mutex sync.Mutex

type PacketSorter struct {
	minHeap []*information.Packet
	length  int
}

func (s PacketSorter) Length() int {
	return s.length - 1
}

func NewInfoSorter() *PacketSorter {
	s := &PacketSorter{}
	s.minHeap = append(s.minHeap, nil)
	s.length++
	return s
}

func (s *PacketSorter) Print() {
	for i := 1; i < s.length; i++ {
		fmt.Print(s.minHeap[i].Timestamp(), " ")
	}
	fmt.Println()
}

func (s *PacketSorter) Append(info *information.Packet) {
	//mutex.Lock()
	s.minHeap = append(s.minHeap, info)
	s.length++
	//mutex.Unlock()
	s.adjustUp(s.length - 1)
}

func (s *PacketSorter) Take() (*information.Packet, error) {
	//mutex.Lock()
	if s.length <= 1 {
		return nil, errors.New("take info from an empty sorter")
	}
	toTake := s.minHeap[1]
	s.minHeap[1], s.minHeap[s.length-1] = s.minHeap[s.length-1], s.minHeap[1]
	s.length--
	s.minHeap = s.minHeap[:s.length]
	//mutex.Unlock()
	s.adjustDown()
	return toTake, nil
}

func (s *PacketSorter) adjustDown() {
	i := 1
	for i <= (s.length-1)/2 {
		// left child default
		min := 2 * i
		// if right child exists and bigger
		if s.length > min+1 && s.minHeap[min+1].Timestamp() < s.minHeap[min].Timestamp() {
			min += 1
		}
		if s.minHeap[i].Timestamp() > s.minHeap[min].Timestamp() {
			s.minHeap[min], s.minHeap[i] = s.minHeap[i], s.minHeap[min] // swap
		}
		i = min
	}
}

func (s *PacketSorter) adjustUp(i int) {
	for ; i >= 2; i = i / 2 {
		if s.minHeap[i].Timestamp() < s.minHeap[i/2].Timestamp() {
			s.minHeap[i], s.minHeap[i/2] = s.minHeap[i/2], s.minHeap[i]
		} else {
			return
		}
	}
}
