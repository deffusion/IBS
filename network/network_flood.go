package network

import (
	"fmt"
	"github.com/deffusion/IBS/node"
	"github.com/deffusion/IBS/node/routing"
	"math/rand"
)

func NewFloodNode(id int, uploadBandwidth int, region string, config map[string]int) node.Node {
	return node.NewBasicNode(
		uint64(id),
		//hash.Hash64(uint64(id)),
		uploadBandwidth,
		id,
		region,
		routing.NewFloodTable(config["tableSize"], config["degree"]),
	)
}

type FloodNet struct {
	tableSize int
	degree    int
	*BaseNetwork
}

func NewFloodNet(size, tableSize, degree int) *FloodNet {
	fmt.Println("degree:", degree)
	// bootNode is used for message generation (from node) only here
	bootNode := node.NewBasicNode(0, 0, 0, "", nil)
	net := NewBasicNetwork(bootNode)
	config := make(map[string]int)
	config["tableSize"] = tableSize
	config["degree"] = degree
	net.generateNodes(size, NewFloodNode, config)
	fNet := &FloodNet{
		tableSize,
		degree,
		net,
	}
	fNet.initConnections()
	return fNet
}

// Introduce : return n nodes
func (fNet *FloodNet) Introduce(n int) []node.Node {
	var nodes []node.Node
	for i := 0; i < n; i++ {
		r := rand.Intn(fNet.Size()) + 1 // zero is the index of msg generator
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
		//cnts = append(cnts, fNet.degree-connectCount)
		peers := fNet.Introduce(fNet.degree - connectCount)
		for _, peer := range peers {
			if fNet.Connect(node, peer, NewBasicPeerInfo) == true {
				//cnt++
			}
		}
		//cnts = append(cnts, cnt)
	}
	//fmt.Println("connect count: ", cnts)
}

func (fNet *FloodNet) introduceAndConnect(n node.Node, f NewPeerInfo) {
	peers := fNet.Introduce(fNet.tableSize)
	for _, peer := range peers {
		fNet.Connect(n, peer, f)
	}
}

func (fNet *FloodNet) churn(crashFrom int, routing func(tableSize, degree int) routing.Table) int {
	for _, n := range fNet.Nodes {
		if n.Running() == false {
			// it can be seen as the crashed nodes leave the network
			n.ResetRoutingTable(routing(fNet.tableSize, fNet.degree))
			n.Run()
			fNet.introduceAndConnect(n, NewBasicPeerInfo)
		}
	}
	return fNet.NodeCrash(crashFrom)
}

func (fNet *FloodNet) Churn(crashFrom int) int {
	return fNet.churn(crashFrom, routing.NewFloodTable)
}
