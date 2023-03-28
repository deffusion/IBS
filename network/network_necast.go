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
const BucketSize = 12

func NewNecastPeerInfo(n node.Node) routing.PeerInfo {
	return routing.NewNecastPeerInfo(n.Id())
}
func NewNecastNode(index int, downloadBandwidth, uploadBandwidth int, region string, k int) node.Node {
	nodeID := hash.Hash64(uint64(index))
	//nodeID := uint64(index)
	return node.NewNeNode(
		nodeID,
		downloadBandwidth,
		uploadBandwidth,
		index,
		region,
		routing.NewNecastTable(nodeID, k, MinFanOut),
	)
}

func NewNecastNet(size int) *NecastNet {
	// bootNode is used for message generation (from node) only here
	bootNode := node.NewBasicNode(BootNodeID, 0, 0, 0, "", routing.NewNecastTable(BootNodeID, BucketSize, MinFanOut))
	net := NewNetwork(bootNode)
	net.generateNodes(size, NewNecastNode, BucketSize)
	nNet := &NecastNet{
		&KadcastNet{
			BucketSize,
			net,
			num_set.NewSet(),
		},
	}
	nNet.initConnections(NewNecastPeerInfo)
	return nNet
}

func (nNet *NecastNet) Churn(crashFrom int) int {
	for _, n := range nNet.Nodes {
		if n.Running() == false {
			n.ResetRoutingTable(routing.NewNecastTable(n.Id(), BucketSize, MinFanOut))
			n.Run()
			nNet.introduceAndConnect(n, NewNecastPeerInfo)
		}
	}
	return nNet.NodeCrash(crashFrom)
}