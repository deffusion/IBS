package main

import (
	"fmt"
	"github.com/deffusion/IBS/network"
	"log"
)

func main() {
	// for i := 0; i < 20; i++ {
	// 	fmt.Println("=====", i, "=====")
	// 	once()
	// }
	once()
}

func once() {
	var NetSize = 1000
	var LogFactor = NetSize
	var NMessage = 10 * NetSize
	//var NMessage = 1
	// var CrashInterval = 60_000_000 // 60s
	var CrashInterval = 3000000_000_000  // large enough
	var PacketGenerationInterval = 2_000 // ms

	var k = 15
	var beta = 1
	//packetStore = make(map[int]*information.BasicPacket)
	//net := network.NewFloodNet(NetSize, 10)
	net := network.NewKadcastNet(NetSize, k, beta)
	// net := network.NewNecastNet(NetSize, k, beta)
	log.Print("net readys")
	fmt.Printf("NetSize: %d, NMessage: %d, PacketGenerationInterval: %d(μs), CrashInterval: %d(μs)\n",
		NetSize, NMessage, PacketGenerationInterval, CrashInterval)
	// cntCrash := net.Churn(1)
	// fmt.Println("first crashed: ", cntCrash)
	// cntInfest := net.Infest(1)
	// fmt.Println("first infested: ", cntInfest)
	sim := New(net, NMessage, LogFactor, CrashInterval, PacketGenerationInterval)
	sim.Run(true)
	sim.Statistic()
}
