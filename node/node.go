package node

import (
	"IBS/node/routing"
	"fmt"
)

type Node struct {
	id                uint64
	region            string
	routingTable      routing.Table
	downloadBandwidth int // byte/s
	uploadBandwidth   int

	//TsLastReceived int64 // the time(μs) when last packet was received
	TsLastSending int64 // the time(μs) when last packet was sent

	receivedPackets map[int]int64 // id -> delay
	running         bool
}

func NewNode(id uint64, downloadBandwidth, uploadBandwidth int, region string, table routing.Table) *Node {
	return &Node{
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

func (n *Node) Id() uint64 {
	return n.id
}
func (n *Node) Region() string {
	return n.region
}
func (n *Node) DownloadBandwidth() int {
	return n.downloadBandwidth
}
func (n *Node) UploadBandwidth() int {
	return n.uploadBandwidth
}

//func (n *Node) LastDelay() int64 {
//	return n.lastDelay
//}
//func (n *Node) SetLastDelay(delay int64) bool {
//	if delay < n.lastDelay {
//		return false
//	}
//	n.lastDelay = delay
//	return true
//}

func (n *Node) NoRoomForNewPeer() bool {
	return n.routingTable.Length() >= n.routingTable.PeerLimit()
}

func (n *Node) RoutingTableLength() int {
	return n.routingTable.Length()
}

func (n *Node) AddPeer(peerInfo routing.PeerInfo) {
	err := n.routingTable.AddPeer(peerInfo)
	if err != nil {
		fmt.Println(err)
	}
}

// return id of peers
func (n *Node) PeersToBroadCast(from *Node) *[]uint64 {
	peerIDs := n.routingTable.PeersToBroadcast(from.id)
	return &peerIDs
}

func (n *Node) Received(msgId int, timestamp int64) bool {
	_, ok := n.receivedPackets[msgId]
	if !ok {
		n.receivedPackets[msgId] = timestamp
	}
	return ok
}
func (n *Node) NumReceivedPackets() int {
	return len(n.receivedPackets)
}

func (n *Node) Running() bool {
	return n.running
}
func (n *Node) PrintTable() {
	n.routingTable.PrintTable()
}

func (n *Node) Run() {
	n.running = true
}
func (n *Node) Stop() {
	n.running = false
}
