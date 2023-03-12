package main

import (
	"IBS/information"
	"IBS/network"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

const NetSize = 10000
const RecordUnit = NetSize / 10
const NMessage = 1

//func NewBasicPeerInfo(n *node.Node) routing.PeerInfo {
//	return routing.NewBasicPeerInfo(n.Id())
//}

func main() {
	net := network.NewFloodNet(NetSize)
	//net := network.NewKadcastNet(NetSize)
	net.NodeCollapse(NetSize / 10)

	var progress []*PacketStatistic

	sorter := NewInfoSorter()
	offset := 0
	for i := 0; i < NMessage; i++ {
		id := net.NodeID(i%NetSize + 1)
		node := net.Node(id)
		// avoid broadcasting from a node is not running
		for node.Running() == false && offset < NetSize-NMessage {
			offset++
			id = net.NodeID((i+offset)%NetSize + 1)
			node = net.Node(id)
		}
		m := information.NewPacket(i, 1<<7, net.BootNode(), node, node, int64(20*i), net.Network)
		ps := NewPacketStatistic()
		ps.Timestamps[0] = m.InfoTimestamp()
		progress = append(progress, ps)
		sorter.Append(m)
	}

	t := Run(sorter, progress)
	cnt := 0
	regionCount := map[string]int{}
	for i := 1; i <= NetSize; i++ {
		id := net.NodeID(i)
		//id := uint64(i)
		//net.Node(id).PrintTable()
		nPackets := net.Node(id).NumReceivedPackets()
		regionCount[net.Node(id).Region()]++
		if nPackets == NMessage {
			cnt++
		}
	}
	fmt.Printf("%d/%d nodes received %d packet in %d μs\n", cnt, NetSize, NMessage, t)
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
	//os.Create("packets.json")
	err = ioutil.WriteFile("output/packets.json", b, 0777)
	if err != nil {
		fmt.Println(err)
	}
}
