package main

import (
	"IBS/information"
	"IBS/network"
	"IBS/node"
	"IBS/node/routing"
	"fmt"
)

const NetSize = 1000000
const MaxDegree = 30
const NMessage = 1

func NewFloodPeerInfo(n *node.Node) routing.PeerInfo {
	return routing.NewFloodPeerInfo(n.Id())
}
func NewFloodTable() routing.Table {
	return routing.NewFloodTable(MaxDegree)
}

func main() {

	net := network.NewNetwork()
	msgGenerator := node.NewNode(0, 0, 0, "", nil)
	// 2<<20 = 1M (Byte/s)
	net.GenerateNodes(NetSize, NewFloodTable)
	net.GenerateConnections(MaxDegree)
	//var t1 routing.Table = routing.NewFloodTable(10)
	//node1 := node.NewNode(1, 1<<22, 1<<19, net.Regions[0], t1)
	//var t2 routing.Table = routing.NewFloodTable(10)
	//node2 := node.NewNode(2, 1<<21, 1<<19, net.Regions[1], t2)
	//var t3 routing.Table = routing.NewFloodTable(10)
	//node3 := node.NewNode(3, 1<<21, 1<<17, net.Regions[2], t3)
	//net.Connect(node1, node2, NewFloodPeerInfo)
	//net.Connect(node1, node3, NewFloodPeerInfo)
	//net.Connect(node2, node3, NewFloodPeerInfo)
	//net.Add(node1)
	//net.Add(node2)
	//net.Add(node3)
	//var x Sayer
	//p := information.NewPacket(1, 1024, msgGenerator, node1, node1, 0, net)

	sorter := NewInfoSorter()

	for i := 1; i <= NMessage; i++ {
		sorter.Append(information.NewPacket(i, 1<<20, msgGenerator, net.Node(i), net.Node(i), 0, net))
	}

	t, tFinish := Run(sorter)
	cnt := 0
	regionCount := map[string]int{}
	for i := 1; i <= NetSize; i++ {
		nPackets := net.Node(i).NumReceivedPackets()
		regionCount[net.Node(i).Region()]++
		if nPackets == NMessage {
			cnt++
		} else {
			fmt.Printf("node%d received %d packets\n", i, nPackets)
		}
	}
	fmt.Printf("%d nodes received %d packet in %d μs\n", cnt, NMessage, t)
	fmt.Printf("stopped at %d μs\n", tFinish)
	fmt.Println(regionCount)
}

func Run(sorter *PacketSorter) (int64, int64) {
	t := int64(0)
	tFinish := int64(0)
	n := 0 // num of packets were broadcast
	for sorter.Length() > 0 {
		p, _ := sorter.Take()
		//p.Print()
		packets := p.NextPackets()
		n++
		if n%10000 == 0 {
			fmt.Println(n)
		}
		if len(*packets) > 0 {
			t = p.Timestamp()
		} else {
			tFinish = p.Timestamp()
		}
		for _, packet := range *packets {
			sorter.Append(packet)
		}
	}
	return t, tFinish
}
