package network

import (
	"IBS/information"
	"IBS/node"
	"IBS/node/routing"
)

const BootNodeID = 0

type NewPeerInfo func(node.Node) routing.PeerInfo

func NewBasicPeerInfo(n node.Node) routing.PeerInfo {
	return routing.NewBasicPeerInfo(n.Id())
}

type Network interface {
	BootNode() node.Node
	Node(id uint64) node.Node
	NodeID(i int) uint64
	Connect(a, b node.Node, f NewPeerInfo) bool
	Add(n node.Node, i int)
	Size() int
	NodeCrash(i int) int
	NodeInfest(i int) int
	NewPacketGeneration(timestamp int64) information.Packet
	succeedingPackets(p *information.BasicPacket, IDs *[]uint64) information.Packets
	PacketReplacement(p *information.BasicPacket) information.Packets
	Churn(crashFrom int) int
	OutputNodes()
}
