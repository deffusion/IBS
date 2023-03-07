package network

import (
	"IBS/node"
	"IBS/node/routing"
	"math/rand"
)

func NewFloodPeerInfo(n *node.Node) routing.PeerInfo {
	return routing.NewFloodPeerInfo(n.Id())
}

type Network struct {
	nodes          map[int]*node.Node
	RegionId       map[string]int
	Regions        []string
	DelayOfRegions *[][]int64
}

func NewNetwork() *Network {
	regions := []string{"cn", "uk", "usa"}
	regionId := make(map[string]int)
	for i, region := range regions {
		regionId[region] = i
	}
	// unit: μs (0.000,001s)
	delayOfRegions := &[][]int64{
		// unit: μs (0.000,001s)
		{10_000, 200_000, 250_000},
		{200_000, 3_000, 100_000},
		{250_000, 100_000, 7_000},
	}
	return &Network{
		map[int]*node.Node{},
		regionId,
		regions,
		delayOfRegions,
	}
}

func (net *Network) GenerateNodes(n int, newTable func() routing.Table) {
	nodeDistribution := []float32{0.1, 0.3, 0.6}
	uploadBandwidth := []int{1 << 19, 1 << 18, 1 << 17}
	downloadBandwidth := []int{1 << 22, 1 << 21, 1 << 21}
	for i := 1; i <= n; i++ {
		regionIndex := 0
		r := rand.Float32()
		acc := float32(0)
		for index, f := range nodeDistribution {
			acc += f
			if r >= acc && r < acc+f {
				regionIndex = index
			}
		}

		net.Add(node.NewNode(
			i,
			downloadBandwidth[regionIndex],
			uploadBandwidth[regionIndex],
			net.Regions[regionIndex],
			newTable(),
		))
	}
}

// n: at most n connection per node
func (net *Network) GenerateConnections(n int) {
	for _, node := range net.nodes {
		connectCount := node.RoutingTableLength()
		for connectCount < n {
			r := rand.Intn(net.Size()) + 1
			//if net.Connect(node, net.nodes[r], NewFloodPeerInfo) == true {
			net.Connect(node, net.nodes[r], NewFloodPeerInfo)
			connectCount++
			//}
		}
	}
}

func (net *Network) Node(id int) *node.Node {
	return net.nodes[id]
}
func (net *Network) Connect(a, b *node.Node, f func(*node.Node) routing.PeerInfo) bool {
	if a.Id() == b.Id() {
		return false
	}
	if a.NoRoomForNewPeer() || b.NoRoomForNewPeer() {
		return false
	}
	//fmt.Printf("connect %d to %d\n", a.Id(), b.Id())
	a.AddPeer(f(b))
	b.AddPeer(f(a))
	return true
}

func (net *Network) Add(n *node.Node) {
	net.nodes[n.Id()] = n
}

func (net *Network) Size() int {
	return len(net.nodes)
}

//type InfoSorter interface {
//	Append(*information.Information)
//	Take() *information.Information
//}
