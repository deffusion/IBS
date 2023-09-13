package cmd

import (
	"fmt"
	"github.com/deffusion/IBS/cmd/flags"
	"github.com/deffusion/IBS/network"
	"github.com/deffusion/IBS/simulator"
	"github.com/urfave/cli/v2"
)

var Flood = &cli.Command{
	Name: "flood",
	Flags: []cli.Flag{
		flags.BroadcastPerNode,
		flags.BroadcastInterval,
		flags.NetSize,
	},
	Action: func(ctx *cli.Context) error {
		netSize := ctx.Int("net_size") / 10
		broadcastPerNode := ctx.Int("broadcast_per_node")
		nMessage := broadcastPerNode * netSize
		broadcastInterval := ctx.Int("broadcast_interval")
		degree := 9
		tableSize := 15
		crashInterval := 3000000_000_000 // 3000000s

		//folder := fmt.Sprintf("flood_%d", time.Now().Unix())
		//os.Mkdir(folder, 0777)
		//logFile, err := os.Create(folder + "/output_log.txt")
		//if err != nil {
		//	return err
		//}

		net := network.NewFloodNet(netSize, tableSize, degree)
		fmt.Println(net.Nodes)
		sim := simulator.New(net, nMessage, netSize, crashInterval, broadcastInterval)
		sim.Run(true)
		outputText := sim.Statistic()
		fmt.Println(outputText)
		return nil
	},
}
