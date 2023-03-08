package main

import (
	"IBS/information"
	"IBS/network"
	"IBS/node"
	"IBS/node/hash"
	"IBS/node/routing"
	"fmt"
)

const NetSize = 10000
const MaxDegree = 20
const NMessage = 300
const K = 3

//func NewBasicPeerInfo(n *node.Node) routing.PeerInfo {
//	return routing.NewBasicPeerInfo(n.Id())
//}
func NewFloodNode(id int64, downloadBandwidth, uploadBandwidth int, region string) *node.Node {
	return node.NewNode(
		uint64(id),
		downloadBandwidth,
		uploadBandwidth,
		region,
		routing.NewFloodTable(MaxDegree),
	)
}

func NewKadcastNode(index int64, downloadBandwidth, uploadBandwidth int, region string) *node.Node {
	nodeID := hash.Hash64(uint64(index))
	//nodeID := uint64(index)
	return node.NewNode(
		nodeID,
		downloadBandwidth,
		uploadBandwidth,
		region,
		routing.NewKadcastTable(nodeID, K),
	)
}

func main() {

	net := network.NewNetwork()
	msgGenerator := node.NewNode(0, 0, 0, "", nil)
	// 2<<20 = 1M (Byte/s)
	//net.GenerateNodes(NetSize, NewFloodNode)
	//net.InitFloodConnections(MaxDegree)
	net.GenerateNodes(NetSize, NewKadcastNode)
	net.InitKademliaConnections()
	//id := net.NodeID(uint64(103))
	//net.Node(id).PrintTable()
	//return
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
		id := net.NodeID(uint64(i))
		sorter.Append(information.NewPacket(i, 1<<7, msgGenerator, net.Node(id), net.Node(id), 0, net))
	}

	t := Run(sorter)
	cnt := 0
	regionCount := map[string]int{}
	for i := 1; i <= NetSize; i++ {
		id := net.NodeID(uint64(i))
		//id := uint64(i)
		//net.Node(id).PrintTable()
		nPackets := net.Node(id).NumReceivedPackets()
		regionCount[net.Node(id).Region()]++
		if nPackets == NMessage {
			cnt++
		} else {
			fmt.Printf("node%d received %d packets\n", i, nPackets)
		}
	}
	fmt.Printf("%d nodes received %d packet in %d μs\n", cnt, NMessage, t)
	fmt.Println(regionCount)
}

func Run(sorter *PacketSorter) int64 {
	t := int64(0)
	hop := 0
	coveredNodes := 0
	tFinish := int64(0)
	n := 0 // num of packets were broadcast
	for sorter.Length() > 0 {
		p, _ := sorter.Take()
		//p.Print()
		packets := p.NextPackets()
		n++
		//if n%10000 == 0 {
		//	fmt.Println(n)
		//}
		if p.Redundancy() == false {
			coveredNodes++
			t = p.Timestamp()
			if p.Hop() > hop {
				hop = p.Hop()
			}

		}
		tFinish = p.Timestamp()
		//if len(*packets) > 0 {
		//	t = p.Timestamp()
		//} else {
		//	tFinish = p.Timestamp()
		//}
		for _, packet := range *packets {
			sorter.Append(packet)
		}
	}
	fmt.Printf("stopped at %d(μs), %d packets total\n", tFinish, n)
	fmt.Println("max hop", hop)
	return t
}
