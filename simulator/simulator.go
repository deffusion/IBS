package main

import (
	"IBS/information"
	"IBS/network"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

const NetSize = 1000
const RecordUnit = NetSize / 10
const NMessage = 1

//func NewBasicPeerInfo(n *node.Node) routing.PeerInfo {
//	return routing.NewBasicPeerInfo(n.Id())
//}

func main() {
	//msgGenerator := node.NewNode(0, 0, 0, "", nil)
	// 2<<20 = 1M (Byte/s)
	//net.GenerateNodes(NetSize, NewFloodNode)
	//net.InitFloodConnections(MaxDegree)
	//net := network.FloodNet{}

	//net := network.NewFloodNet(NetSize)
	net := network.NewKadcastNet(NetSize)
	//net.GenerateNodes(NetSize, NewKadcastNode)
	//net.InitKademliaConnections()

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

	var progress []*PacketStatistic

	sorter := NewInfoSorter()
	for i := 0; i < NMessage; i++ {
		id := net.NodeID(uint64(i%NetSize + 1))
		m := information.NewPacket(i, 1<<7, net.BootNode(), net.Node(id), net.Node(id), int64(20*i), net.Network)
		ps := NewPacketStatistic()
		ps.Timestamps[0] = m.InfoTimestamp()
		progress = append(progress, ps)
		sorter.Append(m)
	}

	t := Run(sorter, progress)
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

func Run(sorter *PacketSorter, progress []*PacketStatistic) int64 {
	t := int64(0)
	tFinish := int64(0)
	n := 0 // num of packets were broadcast
	var outputs []OutputPacket
	for sorter.Length() > 0 {
		p, _ := sorter.Take()
		outputs = append(outputs, *NewOutputPacket(p))
		//p.Print()
		packets := p.NextPackets()
		n++
		//if n%10000 == 0 {
		//	fmt.Println(n)
		//}
		if p.Redundancy() == false {
			ps := progress[p.ID()]
			ps.Received++
			if ps.Received%RecordUnit == 0 {
				ps.Timestamps[ps.Received] = p.Timestamp()
			}
			t = p.Timestamp()
			if ps.MaxHop < p.Hop() {
				ps.MaxHop = p.Hop()
			}

		}
		tFinish = p.Timestamp()
		for _, packet := range *packets {
			sorter.Append(packet)
		}
	}
	writePackets(&outputs)
	fmt.Printf("stopped at %d(μs), %d packets total\n", tFinish, n)
	fmt.Println("progress:")
	//for i, statistic := range progress {
	//	fmt.Printf("packet %d start at %d delay=%d\n",
	//		i, statistic.Timestamps[0], statistic.Delay())
	//}

	return t
}

func writePackets(p *[]OutputPacket) {
	b, err := json.Marshal(*p)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(string(b))
	os.Create("packets.json")
	err = ioutil.WriteFile("packets.json", b, 0777)
	if err != nil {
		fmt.Println(err)
	}
}
