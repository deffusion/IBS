package main

import (
	"fmt"
	"github.com/deffusion/IBS/network"
	"github.com/deffusion/IBS/simulator"
	"log"
	"math"
	"os"
	"strings"
	"time"
)

//func main() {
//	app := cli.App{
//		Name:     "IBS",
//		Commands: cmd.Root,
//	}
//	if err := app.Run(os.Args); err != nil {
//		log.Fatal(err)
//	}
//}

func main() {
	for i := 0; i < 10; i++ {
		fmt.Println("=====", i, "=====")
		once()
	}
	//once5ms()
}

func datetime() string {
	t := strings.Split(time.Now().String(), ".")[0]
	return strings.ReplaceAll(t, ":", ".")
}

func once() {
	netSize := 1000
	broadcastPerNode := 10
	nMessage := broadcastPerNode * netSize
	broadcastInterval := 50_000
	k := 15
	beta := 4
	crashInterval := math.MaxInt
	withNE := false
	malicious := false

	if malicious && crashInterval < math.MaxInt {
		log.Fatalln("bad instruction")
	}

	folder := fmt.Sprintf(
		"kad_coverage_%s(beta=%d,interval=%d)",
		datetime(), beta, broadcastInterval,
	)
	if withNE {
		folder = "ne" + folder
	}
	fmt.Println("output will be write to path:", folder)
	os.Mkdir(folder, 0777)
	logFile, err := os.Create(folder + "/output_log.txt")
	if err != nil {
		log.Fatal(err)
		return
	}
	configText := "===== kademlia =====\n"
	if withNE {
		configText = "===== ne-kademlia =====\n"
	}
	configText += fmt.Sprintf("beta: %d, bucket size: %d\n", beta, k)
	configText += fmt.Sprintf("NetSize: %d, NMessage: %d, BroadcastInterval: %d(μs)\n",
		netSize, nMessage, broadcastInterval)
	var InitNet func(size int, k int, beta int) network.Network
	InitNet = network.NewKadcastNet
	if withNE {
		InitNet = network.NewNecastNet
	}
	net := InitNet(netSize, k, beta)
	if malicious {
		cntInfest := net.Infest(1)
		configText += fmt.Sprintf("infested: %d\n", cntInfest)
	} else {
		cntCrash := net.Churn(1, !(crashInterval < math.MaxInt))
		configText += fmt.Sprintf("first crashed: %d\n", cntCrash)
	}
	fmt.Println(configText)
	sim := simulator.New(net, nMessage, netSize, crashInterval, broadcastInterval)
	sim.InitState()
	sim.Run(false, false)
	outputText := sim.Statistic()
	logFile.Write([]byte(configText))
	logFile.Write([]byte(outputText))
	fmt.Println(outputText)
	sim.OutputNodes(folder)
	sim.OutputCoverage(folder)
}
