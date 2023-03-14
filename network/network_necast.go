package network

import (
	"IBS/network/num_set"
	"IBS/node"
	"IBS/node/hash"
	"IBS/node/routing"
)

type NecastNet struct {
	*KadcastNet
}

const MinFanOut = 1

func NewNecastNode(index int64, downloadBandwidth, uploadBandwidth int, region string, k int) node.Node {
	nodeID := hash.Hash64(uint64(index))
	//nodeID := uint64(index)
	return node.NewNeNode(
		nodeID,
		downloadBandwidth,
		uploadBandwidth,
		region,
		routing.NewNecastTable(nodeID, k, MinFanOut),
	)
}

func NewNecastNet(size int) *NecastNet {
	const BucketSize = 15
	// bootNode is used for message generation (from node) only here
	bootNode := node.NewBasicNode(BootNodeID, 0, 0, "", routing.NewNecastTable(BootNodeID, BucketSize, MinFanOut))
	net := NewNetwork(bootNode)
	net.generateNodes(size, NewNecastNode, BucketSize)
	nNet := &NecastNet{
		&KadcastNet{
			BucketSize,
			net,
			num_set.NewSet(),
		},
	}
	nNet.initConnections()
	return nNet
}
