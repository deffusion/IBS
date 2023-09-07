package main

import (
	"github.com/deffusion/IBS/node"
)

type MessageRec struct {
	From       node.Node
	Received   int
	MaxHop     int
	Timestamps map[int]int64
}

func NewPacketStatistic(from node.Node, timestamp int64) *MessageRec {
	return &MessageRec{
		from,
		0,
		0,
		map[int]int64{0: timestamp},
	}
}

func (ps *MessageRec) Delay(last int) int {
	return int(ps.Timestamps[last] - ps.Timestamps[0])
}
