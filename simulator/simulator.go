package main

import (
	"IBS/information"
	"IBS/network"
	"IBS/node"
	"IBS/output"
	"container/heap"
	"fmt"
	"log"
)

type Simulator struct {
	nMessage          int
	logFactor         int
	crashInterval     int
	broadcastInterval int
	lastCrashAt       int64

	net            network.Network
	sorter         *information.PacketSorter
	progress       []*MessageRec // information
	NodeOutput     output.NodeOutput
	coverageOutput output.PacketCoverageOutput
	delayOutput    output.DelayOutput

	endAt               int64
	sentCnt, confirmCnt int
}

func New(net network.Network, nMessage, logFactor, crashInterval, broadcastInterval int) *Simulator {
	return &Simulator{
		nMessage,
		logFactor,
		crashInterval,
		broadcastInterval,
		0,

		net,
		information.NewInfoSorter(),
		[]*MessageRec{},
		output.NewNodeOutput(),
		output.NewCoverageOutput(),
		output.NewDelayOutput(),

		0,
		0,
		0,
	}
}

func (s *Simulator) InitBroadcast() {
	//m := s.net.NewPacketGeneration(0)
	//heap.Push(s.sorter, m)
	//s.progress = append(s.progress, NewPacketStatistic(m.To(), m.Timestamp()))
	for i := 0; i < s.nMessage; i++ {
		m := s.net.NewPacketGeneration(int64(s.broadcastInterval * (i)))
		heap.Push(s.sorter, m)
		s.progress = append(s.progress, NewPacketStatistic(m.To(), m.Timestamp()))
		//newPacketGeneration(net.BaseNetwork, sorter, &progress, int64(PacketGenerationInterval*(i)))
	}
}

// Run the packet replacement process until no packet remains in the sorter
func (s *Simulator) Run() {
	//tFinish := int64(0)
	//outputPackets := output.NewPacketOutput()
	//var outputs []*output.Packet
	fmt.Println("=======len", s.sorter.Len())

	for s.sorter.Len() > 0 {
		p := heap.Pop(s.sorter).(*information.BasicPacket)
		//sorter.Take()
		//switch p.(type) {
		//case *information.TimerPacket:
		//	peers := setTimer(sorter, p.(*information.TimerPacket))
		//	packets := packetStore[p.ID()].succeedingPackets(&peers)
		//	for _, packet := range packets {
		//		sorter.Append(packet)
		//	}
		//case *information.BasicPacket:
		if p.To().Running() == false {
			continue
		}
		if p.To().Malicious() == true {
			continue
		}
		//packet := p.(*information.BasicPacket)
		var succeedingPackets information.Packets
		switch p.To().(type) {
		case *node.NeNode:
			neNode := p.To().(*node.NeNode)
			// the sender will be added if the receiver's bucket have space for it
			if p.From().Id() != network.BootNodeID {
				neNode.AddPeer(network.NewNecastPeerInfo(p.From()))
			}
			// if the packet is sent by this node
			if p.Origin().Id() == neNode.Id() && p.From().Id() != network.BootNodeID {
				s.confirmCnt++
				neNode.Confirmation(p.From().Id(), p.Relay().Id())
			}
		default:
			if p.From().Id() != network.BootNodeID {
				p.To().AddPeer(network.NewBasicPeerInfo(p.From()))
			}
		}
		succeedingPackets = s.net.PacketReplacement(p)
		// churn the network
		if p.Timestamp()-s.lastCrashAt > int64(s.crashInterval) {
			s.lastCrashAt = p.Timestamp()
			fmt.Println("t:", p.Timestamp(), "crashed:", s.net.Churn(1))
		}
		for _, sp := range succeedingPackets {
			heap.Push(s.sorter, sp)
		}

		s.sentCnt++
		//outputs = append(outputs, output.NewPacket(packet))
		if p.Redundancy() == false {
			s.coverageOutput[p.ID()]++
			ps := s.progress[p.ID()]
			ps.Received++
			if ps.Received%(s.net.Size()/5) == 0 {
				ps.Timestamps[ps.Received] = p.Timestamp()
			}
			s.endAt = p.Timestamp()
			if ps.MaxHop < p.Hop() {
				ps.MaxHop = p.Hop()
			}
		}
		//tFinish = p.Timestamp()
		//}
		//p.Print()
		//outputPackets.Append(packet)

	}
	//output.WritePackets(outputs)
	//outputPackets.WritePackets()

}

func (s *Simulator) Statistic() {
	delayOutput := output.NewDelayOutput()
	fmt.Println("progress:")
	for i, statistic := range s.progress {
		delayOutput.Append(i, statistic.Delay(s.net.Size()), statistic.From.Region())
		if i%s.logFactor != 0 {
			continue
		}
		//fmt.Printf("packet %d start at %d delay=%d\n",
		//	i, statistic.Timestamps[0], statistic.Delay())
		fmt.Printf("packet %d coverage:(%d) \n", i, statistic.Received)
	}
	delayOutput.WriteDelay()

	receivedAll := 0
	receivedCnt := 0
	regionCount := map[string]int{}
	bandwidthCount := map[int]int{}
	for i := 1; i <= s.net.Size(); i++ {
		id := s.net.NodeID(i)
		//id := uint64(i)
		//net.Node(id).PrintTable()
		nPackets := s.net.Node(id).NumReceivedPackets()
		receivedCnt += nPackets
		regionCount[s.net.Node(id).Region()]++
		bandwidthCount[s.net.Node(id).UploadBandwidth()]++
		if nPackets == s.nMessage {
			receivedAll++
		}
	}
	fmt.Printf("(%d/%d) received, %d packets totalSent (%d redundancy confirm packet)\n", receivedCnt, (s.net.Size())*s.nMessage, s.sentCnt, s.confirmCnt)
	fmt.Printf("%d/%d nodes received %d packet in %d μs\n", receivedAll, s.net.Size(), s.nMessage, s.endAt)
	fmt.Println("region distribution:", regionCount)
	fmt.Println("upload bandwidth distribution:", bandwidthCount)
	log.Print("end")
	s.coverageOutput.WriteCoverage()
	s.net.OutputNodes()

	//for i, progress := range s.progress {
	//	fmt.Println(i, *progress)
	//}
}
