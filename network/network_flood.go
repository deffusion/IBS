package network

import (
	"IBS/node"
	"IBS/node/routing"
	"math/rand"
)

const MaxDegree = 4

func NewFloodNode(id int64, downloadBandwidth, uploadBandwidth int, region string) *node.Node {
	return node.NewNode(
		uint64(id),
		downloadBandwidth,
		uploadBandwidth,
		region,
		routing.NewFloodTable(MaxDegree),
	)
}

type FloodNet struct {
	*Network
}

func NewFloodNet(size int64) *FloodNet {
	// bootNode is used for message generation (from node) only here
	bootNode := node.NewNode(0, 0, 0, "", routing.NewFloodTable(MaxDegree))
	net := NewNetwork(bootNode)
	net.generateNodes(size, NewFloodNode)
	fNet := &FloodNet{
		net,
	}
	fNet.initConnections()
	return fNet
}

// Introduce : return n nodes
func (fNet *FloodNet) Introduce(n int) []*node.Node {
	var nodes []*node.Node
	for i := 0; i < n; i++ {
		r := rand.Intn(fNet.Size()) + 1 // zero is the msg generator
		//fmt.Println("r", r)
		nodes = append(nodes, fNet.Node(fNet.NodeID(uint64(r))))
	}
	return nodes
}

func (fNet *FloodNet) initConnections() {
	for _, node := range fNet.nodes {
		//fNet.bootNode.AddPeer(NewBasicPeerInfo(node))
		connectCount := node.RoutingTableLength()
		peers := fNet.Introduce(MaxDegree - connectCount)
		for _, peer := range peers {
			fNet.Connect(node, peer, NewBasicPeerInfo)
		}
	}
}
