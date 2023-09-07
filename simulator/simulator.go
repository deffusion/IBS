package simulator

import (
	"container/heap"
	"fmt"
	"github.com/deffusion/IBS/information"
	"github.com/deffusion/IBS/network"
	"github.com/deffusion/IBS/node"
	"github.com/deffusion/IBS/output"
)

type Simulator struct {
	nMessage          int
	logFactor         int
	crashInterval     int
	broadcastInterval int
	lastCrashAt       int64
	broadcastID       int

	net            network.Network
	sorter         *information.PacketSorter
	progress       []*MessageRec // information
	NodeOutput     output.NodeOutput
	coverageOutput output.PacketCoverageOutput
	latencyOutput  output.LatencyOutput

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
		0,

		net,
		information.NewInfoSorter(),
		[]*MessageRec{},
		output.NewNodeOutput(),
		output.NewCoverageOutput(),
		output.NewLatencyOutput(),

		0,
		0,
		0,
	}
}

func (s *Simulator) initAllBroadcast() {
	//m := s.net.NewPacketGeneration(0)
	//heap.Push(s.sorter, m)
	//s.progress = append(s.progress, NewPacketStatistic(m.To(), m.Timestamp()))
	for s.broadcastID < s.nMessage {
		s.initOneBroadcast()
		//newPacketGeneration(net.BaseNetwork, sorter, &progress, int64(PacketGenerationInterval*(i)))
	}
}

func (s *Simulator) initOneBroadcast() {
	m := s.net.NewPacketGeneration(int64(s.broadcastInterval * (s.broadcastID)))
	s.broadcastID++
	heap.Push(s.sorter, m)
	s.progress = append(s.progress, NewPacketStatistic(m.To(), m.Timestamp()))
}

// Run the packet replacement process until no packet remains in the sorter
func (s *Simulator) Run(initAllBroadcast bool) {
	if initAllBroadcast {
		s.initAllBroadcast()
	} else {
		s.initOneBroadcast()
	}
	//tFinish := int64(0)
	//outputPackets := output.NewPacketOutput()
	//var outputs []*output.Packet
	var p *information.BasicPacket
	var sender node.Node
	//malitrans := 0
	//malicious, total := 0, 0
	for s.sorter.Len() > 0 {
		p = heap.Pop(s.sorter).(*information.BasicPacket)
		sender = p.To()
		//sorter.Take()
		//switch p.(type) {
		//case *information.TimerPacket:
		//	peers := setTimer(sorter, p.(*information.TimerPacket))
		//	packets := packetStore[p.ID()].succeedingPackets(&peers)
		//	for _, packet := range packets {
		//		sorter.Append(packet)
		//	}
		//case *information.BasicPacket:
		if sender.Running() == false {
			continue
		}
		s.sentCnt++
		if sender.Malicious() == true {
			//malitrans++
			continue
		}
		//packet := p.(*information.BasicPacket)
		var succeedingPackets information.Packets
		switch dNode := sender.(type) {
		case *node.NeNode:
			// the sender will be added if the receiver's bucket have space for it
			if p.From().Id() != network.BootNodeID {
				dNode.AddPeer(network.NewNecastPeerInfo(p.From()))
			}
			// if the packet is sent by this node
			if p.Origin().Id() == dNode.Id() && p.From().Id() != network.BootNodeID {
				s.confirmCnt++
				dNode.Confirmation(p.From().Id(), p.Relay().Id())
			}
		default:
			if p.From().Id() != network.BootNodeID {
				sender.AddPeer(network.NewBasicPeerInfo(p.From()))
			}
		}
		succeedingPackets = s.net.PacketReplacement(p)
		//succeedingPackets, m, t := s.net.PacketReplacement(p)
		//malicious += m
		//total += t
		// churn the network
		if p.Timestamp()-s.lastCrashAt > int64(s.crashInterval) {
			s.lastCrashAt = p.Timestamp()
			fmt.Println("t:", p.Timestamp(), "crashed:", s.net.Churn(1))
		}
		if !initAllBroadcast && p.Timestamp() > int64(s.broadcastID*s.broadcastInterval) && s.broadcastID < s.nMessage {
			s.initOneBroadcast()
		}
		for _, sp := range succeedingPackets {
			heap.Push(s.sorter, sp)
		}

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
	//fmt.Printf("%d/%d\n", malicious, total)
	//fmt.Println("malicious transmission", malitrans)
	//output.WritePackets(outputs)
	//outputPackets.WritePackets()

}

func (s *Simulator) Statistic() string {
	outputText := ""
	s.latencyOutput = output.NewLatencyOutput()
	//fmt.Println("progress:")
	for i, statistic := range s.progress {
		s.latencyOutput.Append(i, statistic.Delay(s.net.Size()), statistic.From.Region())
		if i%s.logFactor != 0 {
			continue
		}
		//fmt.Printf("packet %d start at %d delay=%d\n",
		//	i, statistic.Timestamps[0], statistic.Delay())
		//fmt.Printf("packet %d coverage:(%d) \n", i, statistic.Received)
	}
	//s.latencyOutput.WriteLatency()

	receivedAll := 0
	receivedCnt := 0
	regionCount := map[string]int{}
	var n node.Node
	bandwidthCount := map[int]int{}
	for i := 1; i <= s.net.Size(); i++ {
		id := s.net.NodeID(i)
		//id := uint64(i)
		//net.Node(id).PrintTable()
		n = s.net.Node(id)
		nPackets := n.NumReceivedPackets()
		//receivedCnt += nPackets
		regionCount[n.Region()]++
		bandwidthCount[n.UploadBandwidth()]++
		if nPackets == s.nMessage {
			receivedAll++
		}
	}
	for _, cnt := range s.coverageOutput {
		receivedCnt += cnt
	}
	outputText += fmt.Sprintf("%d received, %d packets totalSent (%d redundancy confirm packet)\n", receivedCnt, s.sentCnt, s.confirmCnt)
	outputText += fmt.Sprintf("%d/%d nodes received %d packet in %d μs\n", receivedAll, s.net.Size(), s.nMessage, s.endAt)
	outputText += fmt.Sprintf("region distribution:%v\n", regionCount)
	outputText += fmt.Sprintf("upload bandwidth distribution:%v\n", bandwidthCount)
	return outputText
	//fmt.Printf("%d received, %d packets totalSent (%d redundancy confirm packet)\n", receivedCnt, s.sentCnt, s.confirmCnt)
	//fmt.Printf("%d/%d nodes received %d packet in %d μs\n", receivedAll, s.net.Size(), s.nMessage, s.endAt)
	//fmt.Println("region distribution:", regionCount)
	//fmt.Println("upload bandwidth distribution:", bandwidthCount)
	//log.Print("end")
	//s.coverageOutput.WriteCoverage()
	//s.net.OutputNodes()

	//for i, progress := range s.progress {
	//	fmt.Println(i, *progress)
	//}
}

func (s *Simulator) OutputCoverage(folder string) {
	s.coverageOutput.WriteCoverage(folder)
}
func (s *Simulator) OutputLatency(folder string) {
	s.latencyOutput.WriteLatency(folder)
}
func (s *Simulator) OutputNodes(folder string) {
	s.net.OutputNodes(folder)
}

//func (s Simulator) OutputPackets()  {
//
//}
