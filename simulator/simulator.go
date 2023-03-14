package main

import (
	"IBS/information"
	"IBS/network"
	"IBS/node"
	"IBS/output"
	"fmt"
)

const NetSize = 10000
const RecordUnit = NetSize / 10
const NMessage = 1

//func NewBasicPeerInfo(n *node.BasicNode) routing.PeerInfo {
//	return routing.NewBasicPeerInfo(n.Id())
//}

var packetStore map[int]*information.BasicPacket

func main() {
	packetStore = make(map[int]*information.BasicPacket)
	//net := network.NewFloodNet(NetSize)
	//net := network.NewKadcastNet(NetSize)
	net := network.NewNecastNet(NetSize)

	var outputNodes []*output.Node
	for _, n := range net.Nodes {
		outputNodes = append(outputNodes, output.NewNode(n))
	}
	output.WriteNodes(outputNodes)
	//net.NodeCollapse(NetSize / 10 * 5)

	var progress []*PacketStatistic

	sorter := NewInfoSorter()
	offset := 0
	for i := 0; i < NMessage; i++ {
		id := net.NodeID(i%NetSize + 1)
		n := net.Node(id)
		// avoid broadcasting from a node is not running
		for n.Running() == false && offset < NetSize-NMessage {
			offset++
			id = net.NodeID((i+offset)%NetSize + 1)
			n = net.Node(id)
		}
		switch n.(type) {
		case *node.NeNode:
			neNode := n.(*node.NeNode)
			neNode.NewBroadcastTask(i)
			fmt.Println(neNode.Id(), "tasks", *neNode.Tasks[i])
		}
		m := information.NewBasicPacket(i, 1<<7, n, net.BootNode(), n, nil, int64(20*i), net.Network)
		packetStore[m.ID()] = m
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
	var outputs []*output.Packet
	for sorter.Length() > 0 {
		p, _ := sorter.Take()
		switch p.(type) {
		case *information.TimerPacket:
			peers := setTimer(sorter, p.(*information.TimerPacket))
			packets := packetStore[p.ID()].NextPackets(peers)
			for _, packet := range packets {
				sorter.Append(packet)
			}
		case *information.BasicPacket:
			packet := p.(*information.BasicPacket)
			switch p.To().(type) {
			case *node.NeNode:
				neNode := p.To().(*node.NeNode)
				// if the packet is sent by this node
				_, ok := neNode.Tasks[p.ID()]
				if ok {
					neNode.Confirm(p.ID(), p.From().Id())
				}
			}
			broadcast(sorter, packet)
			n++
			//if n%10000 == 0 {
			//	fmt.Println(n)
			//}
			outputs = append(outputs, output.NewPacket(packet))
			if packet.Redundancy() == false {
				ps := progress[p.ID()]
				ps.Received++
				if ps.Received%RecordUnit == 0 {
					ps.Timestamps[ps.Received] = p.Timestamp()
				}
				t = p.Timestamp()
				if ps.MaxHop < packet.Hop() {
					ps.MaxHop = packet.Hop()
				}
			}
			tFinish = p.Timestamp()
		}
		//p.Print()

	}
	output.WritePackets(outputs)
	fmt.Printf("stopped at %d(μs), %d packets total\n", tFinish, n)
	fmt.Println("progress:")
	//for i, statistic := range progress {
	//	fmt.Printf("packet %d start at %d delay=%d\n",
	//		i, statistic.Timestamps[0], statistic.Delay())
	//}

	return t
}

func setTimer(sorter *PacketSorter, p *information.TimerPacket) *[]uint64 {
	neNode := p.To().(*node.NeNode)
	peerIDs := neNode.PeersFromTask(p.ID(), -1)
	// next timer
	if len(*peerIDs) > 0 {
		sorter.Append(p.NextPacket(10_000)) // 10ms
	}
	return peerIDs
}

func broadcast(sorter *PacketSorter, p *information.BasicPacket) {
	var peers *[]uint64
	peers = p.To().PeersToBroadCast(p.From())
	switch p.To().(type) {
	case *node.NeNode:
		neNode := p.To().(*node.NeNode)
		if p.From().Id() == 0 {
			neNode.NewBroadcastTask(p.ID())
			peers = neNode.PeersFromTask(p.ID(), -1)
		}
		if neNode.Id() != p.Origin().Id() && neNode.IsNeighbour(p.Origin().Id()) {
			*peers = append(*peers, p.Origin().Id())
		}
	}
	//fmt.Println("send to peers", *peers)
	packets := p.NextPackets(peers)
	for _, packet := range packets {
		sorter.Append(packet)
	}
}
