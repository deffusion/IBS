package cmd

import (
	"fmt"
	"github.com/deffusion/IBS/cmd/flags"
	"github.com/deffusion/IBS/network"
	"github.com/deffusion/IBS/simulator"
	"github.com/urfave/cli/v2"
	"os"
	"time"
)

var Kademlia = &cli.Command{
	Name: "kademlia",
	Subcommands: []*cli.Command{
		kadLatency,
	},
}

var kadLatency = &cli.Command{
	Name: "latency",
	Flags: []cli.Flag{
		flags.K,
		flags.Beta,
		flags.BroadcastPerNode,
		flags.BroadcastInterval,
		flags.NetSize,
	},
	Action: func(ctx *cli.Context) error {
		netSize := ctx.Int("net_size")
		broadcastPerNode := ctx.Int("broadcast_per_node")
		nMessage := broadcastPerNode * netSize
		broadcastInterval := ctx.Int("broadcast_interval")
		k := ctx.Int("k")
		beta := ctx.Int("beta")
		var crashInterval = 3000000_000_000 // 3000000s

		folder := fmt.Sprintf(
			"kad_latency_%d(beta=%d,interval=%d)",
			time.Now().Unix(), beta, broadcastInterval,
		)
		fmt.Println("output will be write to path:", folder)
		os.Mkdir(folder, 0777)
		logFile, err := os.Create(folder + "/output_log.txt")
		if err != nil {
			return err
		}
		configText := "===== kademlia =====\n"
		configText += fmt.Sprintf("beta: %d, bucket size: %d\n", beta, k)
		configText += fmt.Sprintf("NetSize: %d, NMessage: %d, BroadcastInterval: %d(μs)\n",
			netSize, nMessage, broadcastInterval)
		logFile.Write([]byte(configText))
		fmt.Println(configText)
		net := network.NewKadcastNet(netSize, k, beta)
		sim := simulator.New(net, nMessage, netSize, crashInterval, broadcastInterval)
		sim.Run(true)
		outputText := sim.Statistic()
		logFile.Write([]byte(outputText))
		fmt.Println(outputText)
		sim.OutputNodes(folder)
		sim.OutputLatency(folder)
		return nil
	},
}
