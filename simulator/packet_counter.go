package main

type PacketStatistic struct {
	Received   int
	MaxHop     int
	Timestamps map[int]int64
}

func NewPacketStatistic() *PacketStatistic {
	return &PacketStatistic{
		0,
		0,
		map[int]int64{},
	}
}

func (ps *PacketStatistic) Delay() int {
	return int(ps.Timestamps[NetSize] - ps.Timestamps[0])
}
