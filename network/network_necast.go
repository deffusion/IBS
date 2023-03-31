package network

import (
	"IBS/network/num_set"
	"IBS/node"
	"IBS/node/hash"
	"IBS/node/routing"
	"fmt"
)

type NecastNet struct {
	*KadcastNet
}

func NewNecastPeerInfo(n node.Node) routing.PeerInfo {
	return routing.NewNecastPeerInfo(n.Id())
}
func NewNecastNode(index int, uploadBandwidth int, region string, k int) node.Node {
	nodeID := hash.Hash64(uint64(index))
	//nodeID := uint64(index)
	return node.NewNeNode(
		nodeID,
		uploadBandwidth,
		index,
		region,
		routing.NewNecastTable(nodeID, k, Beta),
	)
}

func NewNecastNet(size int) *NecastNet {
	fmt.Println("===== ne-kademlia =====")
	fmt.Println("beta:", Beta, "bucket size:", K)
	// bootNode is used for message generation (from node) only here
	bootNode := node.NewBasicNode(BootNodeID, 0, 0, "", routing.NewNecastTable(BootNodeID, K, Beta))
	net := NewNetwork(bootNode)
	net.generateNodes(size, NewNecastNode, K)
	nNet := &NecastNet{
		&KadcastNet{
			K,
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
			n.ResetRoutingTable(routing.NewNecastTable(n.Id(), K, Beta))
			n.Run()
			nNet.introduceAndConnect(n, NewNecastPeerInfo)
		}
	}
	return nNet.NodeCrash(crashFrom)
}
