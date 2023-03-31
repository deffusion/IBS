package network

import (
	"IBS/node"
	"IBS/node/hash"
	"IBS/node/routing"
	"fmt"
	"math/rand"
)

func NewFloodNode(id int, uploadBandwidth int, region string, maxDegree int) node.Node {
	return node.NewBasicNode(
		//uint64(id),
		hash.Hash64(uint64(id)),
		uploadBandwidth,
		id,
		region,
		routing.NewFloodTable(maxDegree),
	)
}

type FloodNet struct {
	MaxDgree int
	*Network
}

func NewFloodNet(size int) *FloodNet {
	maxDegree := 15
	fmt.Println("degree:", maxDegree)
	// bootNode is used for message generation (from node) only here
	bootNode := node.NewBasicNode(0, 0, 0, "", routing.NewFloodTable(maxDegree))
	net := NewNetwork(bootNode)
	net.generateNodes(size, NewFloodNode, maxDegree)
	fNet := &FloodNet{
		maxDegree,
		net,
	}
	fNet.initConnections()
	return fNet
}

// Introduce : return n nodes
func (fNet *FloodNet) Introduce(n int) []node.Node {
	var nodes []node.Node
	for i := 0; i < n; i++ {
		r := rand.Intn(fNet.Size()) + 1 // zero is the msg generator
		//fmt.Println("r", r)
		nodes = append(nodes, fNet.Node(fNet.NodeID(r)))
	}
	return nodes
}

func (fNet *FloodNet) initConnections() {
	//var cnts []int
	for _, node := range fNet.Nodes {
		//cnt := 0
		//fNet.bootNode.AddPeer(NewBasicPeerInfo(node))
		connectCount := node.RoutingTableLength()
		//cnts = append(cnts, fNet.MaxDgree-connectCount)
		peers := fNet.Introduce(fNet.MaxDgree - connectCount)
		for _, peer := range peers {
			if fNet.Connect(node, peer, NewBasicPeerInfo) == true {
				//cnt++
			}
		}
		//cnts = append(cnts, cnt)
	}
	//fmt.Println("connect count: ", cnts)
}
