package node

import "IBS/node/routing"

type Node interface {
	Id() uint64
	Region() string
	DownloadBandwidth() int
	UploadBandwidth() int
	TsLastSending() int64
	SetTsLastSending(int64)
	NoRoomForNewPeer() bool
	RoutingTableLength() int
	AddPeer(routing.PeerInfo)
	PeersToBroadCast(Node) *[]uint64
	Received(int, int64) bool
	NumReceivedPackets() int
	Running() bool
	PrintTable()
	Run()
	Stop()
}
