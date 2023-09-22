package cmd

import (
	"errors"
	"fmt"
	"github.com/deffusion/IBS/cmd/flags"
	"github.com/deffusion/IBS/network"
	"github.com/deffusion/IBS/simulator"
	"github.com/urfave/cli/v2"
	"math"
	"os"
)

var Flood = &cli.Command{
	Name: "flood",
	Subcommands: []*cli.Command{
		floodLatency,
		floodCoverage,
	},
}

var floodLatency = &cli.Command{
	Name: "latency",
	Flags: []cli.Flag{
		flags.TableSize,
		flags.Degree,
		flags.BroadcastPerNode,
		flags.BroadcastInterval,
		flags.CrashInterval,
		flags.NetSize,
		flags.WithNE,
	},
	Action: func(ctx *cli.Context) error {
		tableSize := ctx.Int("table_size")
		degree := ctx.Int("degree")

		folder := fmt.Sprintf(
			"flood_latency_%s(tableSize=%d,degree=%d)",
			datetime(), tableSize, degree,
		)
		return floodProcess(ctx, false, folder)
	},
}

var floodCoverage = &cli.Command{
	Name: "coverage",
	Flags: []cli.Flag{
		flags.TableSize,
		flags.Degree,
		flags.BroadcastPerNode,
		flags.BroadcastInterval,
		flags.CrashInterval,
		flags.NetSize,
		flags.WithNE,
		flags.MaliciousNode,
	},
	Action: func(ctx *cli.Context) error {
		tableSize := ctx.Int("table_size")
		degree := ctx.Int("degree")

		folder := fmt.Sprintf(
			"flood_coverage_%s(tableSize=%d,degree=%d)",
			datetime(), tableSize, degree,
		)
		return floodProcess(ctx, true, folder)
	},
}

func floodProcess(ctx *cli.Context, disturbNet bool, folder string) error {
	netSize := ctx.Int("net_size")
	//netSize := 30
	broadcastPerNode := ctx.Int("broadcast_per_node")
	nMessage := broadcastPerNode * netSize
	//nMessage := netSize * 3
	broadcastInterval := ctx.Int("broadcast_interval")
	tableSize := ctx.Int("table_size")
	degree := ctx.Int("degree")
	crashInterval := ctx.Int("crash_interval")
	withNE := ctx.Bool("with_ne")

	if disturbNet {
		malicious := ctx.Bool("malicious")
		if malicious && crashInterval < math.MaxInt {
			return errors.New("bad instruction")
		}
	}

	if withNE {
		folder = "ne" + folder
	}

	fmt.Println("output will be write to path:", folder)
	os.Mkdir(folder, 0777)
	logFile, err := os.Create(folder + "/output_log.txt")
	if err != nil {
		return err
	}
	configText := "===== flood =====\n"
	if withNE {
		configText = "===== ne-flood =====\n"
	}
	configText += fmt.Sprintf("table size: %d, degree: %d\n", tableSize, degree)
	configText += fmt.Sprintf("NetSize: %d, NMessage: %d, BroadcastInterval: %d(μs)\n",
		netSize, nMessage, broadcastInterval)
	logFile.Write([]byte(configText))
	fmt.Println(configText)
	//net := network.NewNeFloodNet(netSize, tableSize, degree)
	var InitNet func(size int, k int, beta int) network.Network
	InitNet = network.NewFloodNet
	if withNE {
		InitNet = network.NewNeFloodNet
	}
	net := InitNet(netSize, tableSize, degree)
	sim := simulator.New(net, nMessage, netSize, crashInterval, broadcastInterval)
	if disturbNet {
		malicious := ctx.Bool("malicious")
		if malicious {
			cntInfest := net.Infest(1)
			fmt.Println("infested: ", cntInfest)
		} else {
			cntCrash := net.Churn(1)
			fmt.Println("first crashed: ", cntCrash)
		}
	}
	sim.Run(false)
	outputText := sim.Statistic()
	logFile.Write([]byte(outputText))
	fmt.Println(outputText)
	sim.OutputNodes(folder)
	sim.OutputLatency(folder)
	return nil
}
