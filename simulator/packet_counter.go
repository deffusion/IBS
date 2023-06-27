package main

import (
	"IBS/node"
)

type PacketStatistic struct {
	From       node.Node
	Received   int
	MaxHop     int
	Timestamps map[int]int64
}

func NewPacketStatistic(from node.Node, timestamp int64) *PacketStatistic {
	return &PacketStatistic{
		from,
		0,
		0,
		map[int]int64{0: timestamp},
	}
}

func (ps *PacketStatistic) Delay(last int) int {
	return int(ps.Timestamps[last] - ps.Timestamps[0])
}
