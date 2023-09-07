package flags

import "github.com/urfave/cli/v2"

var NetSize = &cli.IntFlag{
	Name:        "net_size",
	Usage:       "specify the number of nodes in the network",
	Value:       1000,
	DefaultText: "1000",
}

var BroadcastPerNode = &cli.IntFlag{
	Name:        "broadcast_per_node",
	Usage:       "specify the number of broadcast initialized by each node",
	Value:       10,
	DefaultText: "10",
}

var BroadcastInterval = &cli.IntFlag{
	Name:        "broadcast_interval",
	Usage:       "unit: μs(0.001ms)",
	Value:       5_000, // 5ms
	DefaultText: "5000",
}

var K = &cli.IntFlag{
	Name:        "k",
	Usage:       "bucket size of kademlia",
	Value:       15,
	DefaultText: "15",
}

var Beta = &cli.IntFlag{
	Name:        "beta",
	Usage:       "the broadcast redundancy factor β",
	Value:       1,
	DefaultText: "1",
}