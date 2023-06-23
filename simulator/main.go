package main

import (
	"IBS/network"
	"fmt"
	"log"
)

const NetSize = 1000
const RecordUnit = NetSize / 10
const NMessage = 10 * NetSize
const LogUnit = NetSize
const PacketGenerationInterval = 2_000 // ms
const CrashFrom = 1

//const CrashSpan = 60_000_000 // s
const CrashSpan = 3000000_000_000 // s
//func NewBasicPeerInfo(n *node.BasicNode) routing.PeerInfo {
//	return routing.NewBasicPeerInfo(n.Id())
//}

//var packetStore map[int]*information.BasicPacket
var lastCrashAt int64 // timestamp of last crash

//var lastPacketGeneratedAt int64

func main() {
	//packetStore = make(map[int]*information.BasicPacket)
	//net := network.NewFloodNet(NetSize)
	net := network.NewKadcastNet(NetSize)
	//net := network.NewNecastNet(NetSize)
	log.Print("net ready")
	fmt.Printf("NetSize: %d, NMessage: %d, PacketGenerationInterval: %d(μs), CrashSpan: %d(μs)\n",
		NetSize, NMessage, PacketGenerationInterval, CrashSpan)
	cntCrash := net.Churn(CrashFrom)
	fmt.Println("first crashed: ", cntCrash)
	//cntInfest := net.Infest(CrashFrom)
	//fmt.Println("first infested: ", cntInfest)
	sim := New(net)
	sim.InitBroadcast()
	sim.Run()
	sim.Statistic()
	//t = _t
	//totalSent += _total

}
