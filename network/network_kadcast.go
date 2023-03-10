package network

import (
	"IBS/node"
	"IBS/node/routing"
	"math/rand"
)

const K = 2

func NewKadcastNode(index int64, downloadBandwidth, uploadBandwidth int, region string) *node.Node {
	//nodeID := hash.Hash64(uint64(index))
	nodeID := uint64(index)
	return node.NewNode(
		nodeID,
		downloadBandwidth,
		uploadBandwidth,
		region,
		routing.NewKadcastTable(nodeID, K),
	)
}

type KadcastNet struct {
	*Network
}

func NewKadcastNet(size int64) *KadcastNet {
	// bootNode is used for message generation (from node) only here
	bootNode := node.NewNode(BootNodeID, 0, 0, "", routing.NewKadcastTable(BootNodeID, K))
	net := NewNetwork(bootNode)
	net.generateNodes(size, NewKadcastNode)
	kNet := &KadcastNet{
		net,
	}
	kNet.initConnections()
	return kNet
}

// Introduce : return n nodes
func (kNet *KadcastNet) Introduce(n int) []*node.Node {
	var nodes []*node.Node
	for i := 0; i < n; i++ {
		r := rand.Intn(kNet.Size()) + 1 // zero is the msg generator
		//fmt.Println("r", r)
		nodes = append(nodes, kNet.Node(kNet.NodeID(uint64(r))))
	}
	return nodes
}

func (kNet *KadcastNet) initConnections() {
	for _, node := range kNet.nodes {
		//kNet.bootNode.AddPeer(NewBasicPeerInfo(node))
		connectCount := node.RoutingTableLength()
		peers := kNet.Introduce(MaxDegree - connectCount)
		for _, peer := range peers {
			kNet.Connect(node, peer, NewBasicPeerInfo)
		}
	}
}
