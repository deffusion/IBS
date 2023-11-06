package main

import (
	"github.com/deffusion/IBS/cmd"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	app := cli.App{
		Name:     "IBS",
		Commands: cmd.Root,
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

//func main() {
//	//for i := 0; i < 10; i++ {
//	//	fmt.Println("=====", i, "=====")
//	//	once()
//	//}
//	once()
//}
//
//func datetime() string {
//	t := strings.Split(time.Now().String(), ".")[0]
//	return strings.ReplaceAll(t, ":", ".")
//}
//
//func once() {
//	netSize := 1000
//	broadcastPerNode := 10
//	nMessage := broadcastPerNode * netSize
//	broadcastInterval := 50_000
//	k := 31
//	beta := 6
//	crashInterval := math.MaxInt
//	withNE := true
//	malicious := true
//
//	if malicious && crashInterval < math.MaxInt {
//		log.Fatalln("bad instruction")
//	}
//
//	folder := fmt.Sprintf(
//		"flood_coverage_%s(beta=%d,interval=%d)",
//		datetime(), beta, broadcastInterval,
//	)
//	if withNE {
//		folder = "ne" + folder
//	}
//	configText := "===== flood =====\n"
//	if withNE {
//		configText = "===== ne-flood =====\n"
//	}
//	configText += fmt.Sprintf("beta: %d, bucket size: %d\n", beta, k)
//	configText += fmt.Sprintf("NetSize: %d, NMessage: %d, BroadcastInterval: %d(μs)\n",
//		netSize, nMessage, broadcastInterval)
//	var InitNet func(int, int, int) network.Network
//	InitNet = network.NewFloodNet
//	if withNE {
//		InitNet = network.NewNeFloodNet
//	}
//	net := InitNet(netSize, k, beta)
//	if malicious {
//		cntInfest := net.Infest(1)
//		folder = folder + "_mal"
//		configText += fmt.Sprintf("infested: %d\n", cntInfest)
//	} else {
//		cntCrash := net.Churn(1, !(crashInterval < math.MaxInt))
//		configText += fmt.Sprintf("first crashed: %d\n", cntCrash)
//	}
//	fmt.Println(configText)
//	sim := simulator.New(net, nMessage, netSize, crashInterval, broadcastInterval)
//	sim.InitState()
//	crashText := sim.Run(false, false)
//	os.Mkdir(folder, 0777)
//	fmt.Println("output will be write to path:", folder)
//	logFile, err := os.Create(folder + "/output_log.txt")
//	if err != nil {
//		log.Fatal(err)
//		return
//	}
//	logFile.Write([]byte(configText))
//	logFile.Write([]byte(crashText))
//	outputText := sim.Statistic()
//	logFile.Write([]byte(outputText))
//	fmt.Println(outputText)
//	sim.OutputNodes(folder)
//	sim.OutputCoverage(folder)
//	sim.OutputReceived(folder)
//}
