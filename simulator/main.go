package main

import (
	"IBS/network"
	"fmt"
	"log"
)

func main() {
	var NetSize = 1000
	var LogFactor = NetSize
	var NMessage = 10 * NetSize
	var CrashInterval = 3000000_000_000  // s
	var PacketGenerationInterval = 2_000 // ms
	//packetStore = make(map[int]*information.BasicPacket)
	//net := network.NewFloodNet(NetSize, 10)
	//net := network.NewKadcastNet(NetSize, 15, 2)
	net := network.NewNecastNet(NetSize, 15, 2)
	log.Print("net readys")
	fmt.Printf("NetSize: %d, NMessage: %d, PacketGenerationInterval: %d(μs), CrashInterval: %d(μs)\n",
		NetSize, NMessage, PacketGenerationInterval, CrashInterval)
	cntCrash := net.Churn(1)
	fmt.Println("first crashed: ", cntCrash)
	//cntInfest := net.Infest(CrashFrom)
	//fmt.Println("first infested: ", cntInfest)
	sim := New(net, NMessage, LogFactor, CrashInterval, PacketGenerationInterval)
	sim.InitBroadcast()
	sim.Run()
	sim.Statistic()
}
