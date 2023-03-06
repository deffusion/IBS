package main

import (
	"IBS/information"
	"IBS/network"
	"IBS/node"
	"IBS/node/routing"
)

func NewFloodTable() routing.Table {
	return routing.NewFloodTable(30)
}

func main() {

	net := network.NewNetwork()
	msgGenerator := node.NewNode(0, 0, 0, "", nil)
	// 2<<20 = 1M (Byte/s)
	net.GenerateNodes(10000, NewFloodTable)
	net.GenerateConnections(30)
	//var t1 routing.Table = routing.NewFloodTable(10)
	//node1 := node.NewNode(1, 1<<22, 1<<19, "cn", t1)
	//var t2 routing.Table = routing.NewFloodTable(10)
	//node2 := node.NewNode(2, 1<<21, 1<<18, "uk", t2)
	//var t3 routing.Table = routing.NewFloodTable(10)
	//node3 := node.NewNode(3, 1<<21, 1<<17, "usa", t3)
	//net.Connect(node1, node2, NewFloodPeerInfo)
	//net.Connect(node1, node3, NewFloodPeerInfo)
	//net.Connect(node2, node3, NewFloodPeerInfo)
	//var x Sayer
	p := information.NewPacket(1, 1024, msgGenerator, net.Node(1), net.Node(1), 0, net)

	sorter := NewInfoSorter()
	sorter.Append(p)

	Run(sorter, net)
}

func Run(sorter *PacketSorter, net *network.Network) {
	for sorter.Length() > 0 {
		p, _ := sorter.Take()
		//p.Print()
		packets := p.NextPackets()
		for _, packet := range *packets {
			sorter.Append(packet)
		}
	}
}
