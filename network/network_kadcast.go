package network

import (
	"IBS/network/num_set"
	"IBS/node"
	"IBS/node/hash"
	"IBS/node/routing"
	"log"
)

func NewKadcastNode(index int, downloadBandwidth, uploadBandwidth int, region string, k int) node.Node {
	nodeID := hash.Hash64(uint64(index))
	//nodeID := uint64(index)
	return node.NewBasicNode(
		nodeID,
		downloadBandwidth,
		uploadBandwidth,
		index,
		region,
		routing.NewKadcastTable(nodeID, k),
	)
}

type KadcastNet struct {
	K int
	*Network
	idSet *num_set.Set
}

func NewKadcastNet(size int) *KadcastNet {
	const K = 1
	// bootNode is used for message generation (from node) only here
	bootNode := node.NewBasicNode(BootNodeID, 0, 0, 0, "", routing.NewKadcastTable(BootNodeID, K))
	net := NewNetwork(bootNode)
	net.generateNodes(size, NewKadcastNode, K)
	kNet := &KadcastNet{
		K,
		net,
		num_set.NewSet(),
	}
	kNet.initConnections(NewBasicPeerInfo)
	return kNet
}

// Introduce : return n nodes
func (kNet *KadcastNet) Introduce(id uint64, n int) []node.Node {
	var nodes []node.Node
	for b := 0; b < routing.KeySpaceBits; b++ {
		fakeID, err := routing.FakeIDForBucket(id, b)

		if err != nil {
			log.Fatal(err)
		}
		// 2*n+1 make sures the right keys for the bucket are covered
		// a crude node discovery implementation
		ids := kNet.idSet.Around(fakeID, 2*n+1)
		for _, u := range ids {
			if kNet.Node(u).Running() == true {
				nodes = append(nodes, kNet.Node(u))
			}
		}
		//fmt.Println("introduce to:", id, "peerIDS:", ids)
	}
	return nodes
}

func (kNet *KadcastNet) initConnections(f NewPeerInfo) {
	for _, n := range kNet.Nodes {
		kNet.idSet.Insert(n.Id())
	}
	for _, n := range kNet.Nodes {
		kNet.introduceAndConnect(n, f)
	}
}

func (kNet *KadcastNet) introduceAndConnect(n node.Node, f NewPeerInfo) {
	peers := kNet.Introduce(n.Id(), kNet.K)
	for _, peer := range peers {
		kNet.Connect(n, peer, f)
	}
}

func (kNet *KadcastNet) Churn(crashFrom int) int {
	for _, n := range kNet.Nodes {
		if n.Running() == false {
			// it can be seen as the crashed nodes leave the network
			// and some new nodes entered
			n.ResetRoutingTable(routing.NewKadcastTable(n.Id(), kNet.K))
			n.Run()
			kNet.introduceAndConnect(n, NewNecastPeerInfo)
		}
	}
	return kNet.NodeCrash(crashFrom)
}
