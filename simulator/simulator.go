package main

import (
	"IBS/information"
	"IBS/network"
	"IBS/node"
	"fmt"
	"log"
)

const NetSize = 100
const Collapse = 50
const RecordUnit = NetSize / 10
const NMessage = 10_000

//func NewBasicPeerInfo(n *node.BasicNode) routing.PeerInfo {
//	return routing.NewBasicPeerInfo(n.Id())
//}

var packetStore map[int]*information.BasicPacket

func main() {
	log.Print("start")
	packetStore = make(map[int]*information.BasicPacket)
	//net := network.NewFloodNet(NetSize)
	//net := network.NewKadcastNet(NetSize)
	net := network.NewNecastNet(NetSize)
	log.Print("net ready")

	//var outputNodes []*output.Node
	//for _, n := range net.Nodes {
	//	//n.PrintTable()
	//	outputNodes = append(outputNodes, output.NewNode(n))
	//}
	//output.WriteNodes(outputNodes)
	net.NodeCollapse(Collapse)

	var progress []*PacketStatistic

	sorter := NewInfoSorter()
	offset := 0
	t := int64(0)
	total := 0
	totalReceived := 0
	for i := 0; i < NMessage; i++ {
		id := net.NodeID((i+offset)%NetSize + 1)
		n := net.Node(id)
		// avoid broadcasting from a node is not running
		for n.Running() == false {
			offset++
			id = net.NodeID((i+offset)%NetSize + 1)
			n = net.Node(id)
		}
		//switch n.(type) {
		//case *node.NeNode:
		//neNode := n.(*node.NeNode)
		//neNode.NewBroadcastTask(i)
		//log.Println(neNode.Id(), "tasks", *neNode.Tasks[i])
		//}
		m := information.NewBasicPacket(i, 1<<7, n, net.BootNode(), n, nil, int64(1000000*i), net.Network)
		packetStore[m.ID()] = m
		ps := NewPacketStatistic()
		ps.Timestamps[0] = m.InfoTimestamp()
		progress = append(progress, ps)
		sorter.Append(m)
		_t, _total := Run(sorter, progress)
		t = _t
		total += _total
	}
	fmt.Println("progress:")
	for i, statistic := range progress {
		if i%NetSize != 0 {
			continue
		}
		//unit := NetSize / 5
		//fmt.Printf("packet %d start at %d delay=%d\n",
		//	i, statistic.Timestamps[0], statistic.Delay())
		//fmt.Printf("packet %d delay=%d startAt:%d \t", i, statistic.Delay(), statistic.Timestamps[0])
		//fmt.Println(
		//	statistic.Timestamps[unit]-statistic.Timestamps[0],
		//	statistic.Timestamps[2*unit]-statistic.Timestamps[unit],
		//	statistic.Timestamps[3*unit]-statistic.Timestamps[2*unit],
		//	statistic.Timestamps[4*unit]-statistic.Timestamps[3*unit],
		//	statistic.Timestamps[5*unit]-statistic.Timestamps[4*unit])
		fmt.Printf("packet %d coverage:(%d/%d) \n", i, statistic.Received, NetSize-Collapse)
	}

	cnt := 0
	regionCount := map[string]int{}
	for i := 1; i <= NetSize; i++ {
		id := net.NodeID(i)
		//id := uint64(i)
		//net.Node(id).PrintTable()
		nPackets := net.Node(id).NumReceivedPackets()
		totalReceived += nPackets
		regionCount[net.Node(id).Region()]++
		if nPackets == NMessage {
			cnt++
		}
	}
	fmt.Printf("(%d/%d) received, %d packets total\n", totalReceived, (NetSize-Collapse)*NMessage, total)
	fmt.Printf("%d/%d nodes received %d packet in %d μs\n", cnt, NetSize, NMessage, t)
	fmt.Println(regionCount)
}

func Run(sorter *PacketSorter, progress []*PacketStatistic) (int64, int) {
	t := int64(0)
	//tFinish := int64(0)
	n := 0 // num of packets were broadcast
	//var outputs []*output.Packet
	for sorter.Length() > 0 {
		p, _ := sorter.Take()
		switch p.(type) {
		case *information.TimerPacket:
			peers := setTimer(sorter, p.(*information.TimerPacket))
			packets := packetStore[p.ID()].NextPackets(&peers)
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
				if ok && packet.From().Id() != network.BootNodeID {
					//fmt.Println("123", packet.ID(), packet.From().Id(), packet.Relay())
					neNode.Confirm(packet.ID(), packet.From().Id(), packet.Relay().Id())
				}
			}
			broadcast(sorter, packet)
			n++
			//if n%10000 == 0 {
			//	fmt.Println(n)
			//}
			//outputs = append(outputs, output.NewPacket(packet))
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
			//tFinish = p.Timestamp()
		}
		//p.Print()

	}
	//output.WritePackets(outputs)

	return t, n
}

func setTimer(sorter *PacketSorter, p *information.TimerPacket) []uint64 {
	neNode := p.To().(*node.NeNode)
	peerIDs := neNode.PeersFromTask(p.ID(), -1)
	// next timer
	if len(peerIDs) > 0 {
		sorter.Append(p.NextPacket(10_000)) // 10ms
	}
	return peerIDs
}

func broadcast(sorter *PacketSorter, p *information.BasicPacket) {
	var packets information.Packets
	var peers []uint64
	switch p.To().(type) {
	case *node.NeNode:
		neNode := p.To().(*node.NeNode)
		if p.From().Id() == 0 {
			neNode.NewBroadcastTask(p.ID())
			//neNode.PrintTable()
			peers = neNode.PeersFromTask(p.ID(), -1)
		} else {
			peers = *(neNode.PeersToBroadCast(p.From()))
		}
		if neNode.Id() != p.Origin().Id() && neNode.IsNeighbour(p.Origin().Id()) && neNode.Id() != p.Relay().Id() {
			packets = append(packets, p.ConfirmPacket())
		}
	default:
		peers = *(p.To().PeersToBroadCast(p.From()))
	}
	if p.From().Id() != 0 {
		p.To().SetLastSeen(p.From().Id(), p.Timestamp())
		//if err != nil {
		//	fmt.Printf("%d->%d", p.From().Id(), p.To().Id())
		//	fmt.Println(err)
		//}
	}
	//fmt.Println(p.To().Id(), "peers", peers)
	//fmt.Println("send to peers", *peers)
	packets = append(p.NextPackets(&peers), packets...)
	for _, packet := range packets {
		sorter.Append(packet)
	}
}
