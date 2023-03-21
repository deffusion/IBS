package main

import (
	"IBS/information"
	"IBS/network"
	"IBS/node"
	"IBS/output"
	"fmt"
	"log"
)

const NetSize = 1000
const RecordUnit = NetSize / 10
const NMessage = 10 * NetSize
const LogUnit = 100
const PacketGenerationSpan = 2_000 // 10ms
const CrashFrom = 1
const CrashSpan = 30_000_000 // 30s
//func NewBasicPeerInfo(n *node.BasicNode) routing.PeerInfo {
//	return routing.NewBasicPeerInfo(n.Id())
//}

//var packetStore map[int]*information.BasicPacket
var lastCrashAt int64 // timestamp of last crash
var lastPacketIndex int
var lastOriginNodeIndex int
var lastPacketGeneratedAt int64

func main() {
	log.Print("start")
	//packetStore = make(map[int]*information.BasicPacket)
	//net := network.NewFloodNet(NetSize)
	net := network.NewKadcastNet(NetSize)
	//net := network.NewNecastNet(NetSize)
	log.Print("net ready")
	outputNodes := output.NewNodeOutput()
	for _, n := range net.Nodes {
		//n.PrintTable()
		outputNodes.Append(n)
	}
	outputNodes.WriteNodes()
	//cntCrash := net.Churn(CrashFrom)
	//fmt.Println("first crashed: ", cntCrash)

	var progress []*PacketStatistic
	packetCoverage := output.NewCoverageOutput()

	sorter := NewInfoSorter()
	//newPacketGeneration(net.Network, sorter, &progress, 0)
	for i := 0; i < NMessage; i++ {
		newPacketGeneration(net.Network, sorter, &progress, int64(PacketGenerationSpan*i))
	}
	t := int64(0)
	totalSent := 0 // infer bandwidth usage
	totalReceived := 0
	t, totalSent = Run(net, sorter, progress, packetCoverage)
	//t = _t
	//totalSent += _total
	delayOutput := output.NewDelayOutput()
	fmt.Println("progress:")
	for i, statistic := range progress {
		(*delayOutput)[i] = statistic.Delay()
		if i%LogUnit != 0 {
			continue
		}
		unit := NetSize / 5
		//fmt.Printf("packet %d start at %d delay=%d\n",
		//	i, statistic.Timestamps[0], statistic.Delay())
		fmt.Printf("packet %d delay=%d startAt:%d \t", i, statistic.Delay(), statistic.Timestamps[0])
		fmt.Println(
			statistic.Timestamps[unit]-statistic.Timestamps[0],
			statistic.Timestamps[2*unit]-statistic.Timestamps[unit],
			statistic.Timestamps[3*unit]-statistic.Timestamps[2*unit],
			statistic.Timestamps[4*unit]-statistic.Timestamps[3*unit],
			statistic.Timestamps[5*unit]-statistic.Timestamps[4*unit])
		//fmt.Printf("packet %d coverage:(%d/%d) \n", i, statistic.Received, NetSize/2)
	}
	delayOutput.WriteDelay()

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
	fmt.Printf("(%d/%d) received, %d packets totalSent\n", totalReceived, (NetSize)*NMessage, totalSent)
	fmt.Printf("%d/%d nodes received %d packet in %d μs\n", cnt, NetSize, NMessage, t)
	fmt.Println(regionCount)
	log.Print("end")
	packetCoverage.WriteCoverage()
}

func newPacketGeneration(net *network.Network, sorter *PacketSorter, progress *[]*PacketStatistic, timestamp int64) {
	var origin node.Node
	for i := 0; i <= NetSize; i++ {
		lastOriginNodeIndex = (lastOriginNodeIndex)%NetSize + 1
		origin = net.Node(net.NodeID(lastOriginNodeIndex))
		if origin.Running() == true {
			break
		}
	}
	//fmt.Println("packet:", lastPacketIndex, "broadcast from:", lastOriginNodeIndex)
	m := information.NewBasicPacket(lastPacketIndex, 1<<7, origin, net.BootNode(), origin, nil, timestamp, net)
	lastPacketGeneratedAt = timestamp
	lastPacketIndex++
	ps := NewPacketStatistic()
	ps.Timestamps[0] = m.InfoTimestamp()
	*progress = append(*progress, ps)
	sorter.Append(m)
}

func Run(net interface{}, sorter *PacketSorter, progress []*PacketStatistic, packetCoverage *output.PacketCoverageOutput) (int64, int) {
	t := int64(0)
	//tFinish := int64(0)
	n := 0 // num of packets were broadcast
	//var outputs []*output.Packet
	for sorter.Length() > 0 {
		p, _ := sorter.Take()
		//switch p.(type) {
		//case *information.TimerPacket:
		//	peers := setTimer(sorter, p.(*information.TimerPacket))
		//	packets := packetStore[p.ID()].NextPackets(&peers)
		//	for _, packet := range packets {
		//		sorter.Append(packet)
		//	}
		//case *information.BasicPacket:
		if p.To().Running() == false {
			continue
		}
		packet := p.(*information.BasicPacket)
		switch p.To().(type) {
		case *node.NeNode:
			neNode := p.To().(*node.NeNode)
			// the sender will be added if the receiver's bucket have space for it
			if p.From().Id() != network.BootNodeID {
				neNode.AddPeer(network.NewNecastPeerInfo(p.From()))
			}
			// if the packet is sent by this node
			if packet.Origin().Id() == neNode.Id() && packet.From().Id() != network.BootNodeID {
				neNode.Confirmation(packet.From().Id(), packet.Relay().Id())
			}
		default:
			if p.From().Id() != network.BootNodeID {
				p.To().AddPeer(network.NewBasicPeerInfo(p.From()))
			}
		}
		// churn the network
		switch net.(type) {
		case *network.NecastNet:
			//fmt.Println("necast")
			nNet := net.(*network.NecastNet)
			broadcast(nNet.Network, sorter, packet)
			//if packet.Timestamp()-lastCrashAt > CrashSpan {
			//	lastCrashAt = packet.Timestamp()
			//	fmt.Println("t:", packet.Timestamp(), "crashed:", nNet.Churn(CrashFrom))
			//}
			//if p.Timestamp()-lastPacketGeneratedAt > PacketGenerationSpan && lastPacketIndex < NMessage {
			//	newPacketGeneration(nNet.Network, sorter, &progress, p.Timestamp())
			//}
		case *network.KadcastNet:
			//fmt.Println("kadcast")
			kNet := net.(*network.KadcastNet)
			broadcast(kNet.Network, sorter, packet)
			//if packet.Timestamp()-lastCrashAt > CrashSpan {
			//	lastCrashAt = packet.Timestamp()
			//	fmt.Println("t:", packet.Timestamp(), "crashed:", kNet.Churn(CrashFrom))
			//}
			//if p.Timestamp()-lastPacketGeneratedAt > PacketGenerationSpan && lastPacketIndex < NMessage {
			//	newPacketGeneration(kNet.Network, sorter, &progress, p.Timestamp())
			//}
		case *network.FloodNet:
			broadcast(net.(*network.FloodNet).Network, sorter, packet)
			//	TODO
		}

		n++
		//outputs = append(outputs, output.NewPacket(packet))
		if packet.Redundancy() == false {
			(*packetCoverage)[packet.ID()]++
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
		//}
		//p.Print()

	}
	//output.WritePackets(outputs)

	return t, n
}

//func setTimer(sorter *PacketSorter, p *information.TimerPacket) []uint64 {
//	neNode := p.To().(*node.NeNode)
//	peerIDs := neNode.PeersFromTask(p.ID(), -1)
//	// next timer
//	if len(peerIDs) > 0 {
//		sorter.Append(p.NextPacket(10_000)) // 10ms
//	}
//	return peerIDs
//}

func broadcast(net *network.Network, sorter *PacketSorter, p *information.BasicPacket) {
	var packets information.Packets
	var peers []uint64
	switch p.To().(type) {
	case *node.NeNode:
		neNode := p.To().(*node.NeNode)
		peers = *(neNode.PeersToBroadCast(p.From()))
		if neNode.Id() != p.Origin().Id() && neNode.IsNeighbour(p.Origin().Id()) && neNode.Id() != p.Relay().Id() {
			packets = append(packets, p.ConfirmPacket())
		}
	default:
		peers = *(p.To().PeersToBroadCast(p.From()))
	}
	if p.From().Id() != 0 {
		p.To().SetLastSeen(p.From().Id(), p.Timestamp())
	}

	crashCnt := 0
	for i, peerID := range peers {
		peers[i-crashCnt] = peers[i]
		//fmt.Println("peerID", peerID)
		if net.Node(peerID).Running() == false {
			p.To().RemovePeer(peerID)
			crashCnt++
		}
	}
	peers = peers[:len(peers)-crashCnt]

	packets = append(p.NextPackets(&peers), packets...)
	for _, packet := range packets {
		sorter.Append(packet)
	}
}
