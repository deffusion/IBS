package main

import "IBS/node"

type PacketStatistic struct {
	From       node.Node
	Received   int
	MaxHop     int
	Timestamps map[int]int64
}

func NewPacketStatistic(from node.Node) *PacketStatistic {
	return &PacketStatistic{
		from,
		0,
		0,
		map[int]int64{},
	}
}

func (ps *PacketStatistic) Delay() int {
	return int(ps.Timestamps[NetSize] - ps.Timestamps[0])
}
