package node

import (
	"IBS/node/routing"
	"fmt"
)

type BasicNode struct {
	id                uint64
	region            string
	routingTable      routing.Table
	downloadBandwidth int // byte/s
	uploadBandwidth   int

	//TsLastReceived int64 // the time(μs) when last packet was received
	tsLastSending int64 // the time(μs) when last packet was sent

	receivedPackets map[int]int64 // id -> delay
	running         bool
}

func NewBasicNode(id uint64, downloadBandwidth, uploadBandwidth int, region string, table routing.Table) *BasicNode {
	return &BasicNode{
		id,
		region,
		table,
		downloadBandwidth,
		uploadBandwidth,
		//0,
		0,
		map[int]int64{},
		true,
	}
}

func (n *BasicNode) Id() uint64 {
	return n.id
}
func (n *BasicNode) Region() string {
	return n.region
}
func (n *BasicNode) DownloadBandwidth() int {
	return n.downloadBandwidth
}
func (n *BasicNode) UploadBandwidth() int {
	return n.uploadBandwidth
}
func (n *BasicNode) TsLastSending() int64 {
	return n.tsLastSending
}
func (n *BasicNode) SetTsLastSending(t int64) {
	n.tsLastSending = t
}

//func (n *BasicNode) LastDelay() int64 {
//	return n.lastDelay
//}
//func (n *BasicNode) SetLastDelay(delay int64) bool {
//	if delay < n.lastDelay {
//		return false
//	}
//	n.lastDelay = delay
//	return true
//}

func (n *BasicNode) NoRoomForNewPeer() bool {
	return n.routingTable.Length() >= n.routingTable.PeerLimit()
}

func (n *BasicNode) RoutingTableLength() int {
	return n.routingTable.Length()
}

func (n *BasicNode) AddPeer(peerInfo routing.PeerInfo) {
	err := n.routingTable.AddPeer(peerInfo)
	if err != nil {
		fmt.Println(err)
	}
}

// return id of peers
func (n *BasicNode) PeersToBroadCast(from Node) *[]uint64 {
	peerIDs := n.routingTable.PeersToBroadcast(from.Id())
	return &peerIDs
}

func (n *BasicNode) Received(msgId int, timestamp int64) bool {
	_, ok := n.receivedPackets[msgId]
	if !ok {
		n.receivedPackets[msgId] = timestamp
	}
	return ok
}
func (n *BasicNode) NumReceivedPackets() int {
	return len(n.receivedPackets)
}

func (n *BasicNode) Running() bool {
	return n.running
}
func (n *BasicNode) PrintTable() {
	n.routingTable.PrintTable()
}

func (n *BasicNode) Run() {
	n.running = true
}
func (n *BasicNode) Stop() {
	n.running = false
}
