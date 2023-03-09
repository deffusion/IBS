package network

import (
	"IBS/node"
	"IBS/node/routing"
	"math/rand"
)

func NewBasicPeerInfo(n *node.Node) routing.PeerInfo {
	return routing.NewBasicPeerInfo(n.Id())
}
func NewKadcastPeerInfo(n *node.Node) routing.PeerInfo {
	return routing.NewBasicPeerInfo(n.Id())
}

type Network struct {
	nodes          map[uint64]*node.Node
	indexes        map[uint64]uint64
	RegionId       map[string]int
	Regions        []string
	DelayOfRegions *[][]int32
}

func NewNetwork() *Network {
	regions := []string{"cn", "uk", "usa"}
	regionId := make(map[string]int)
	for i, region := range regions {
		regionId[region] = i
	}
	// unit: μs (0.000,001s)
	delayOfRegions := &[][]int32{
		// unit: μs (0.000,001s)
		{10_000, 200_000, 250_000},
		{200_000, 3_000, 100_000},
		{250_000, 100_000, 7_000},
	}
	return &Network{
		map[uint64]*node.Node{},
		map[uint64]uint64{},
		regionId,
		regions,
		delayOfRegions,
	}
}

func (net *Network) GenerateNodes(n int64, newNode func(int64, int, int, string) *node.Node) {
	nodeDistribution := []float32{0.3, 0.1, 0.6}
	uploadBandwidth := []int{1 << 19, 1 << 18, 1 << 17}
	downloadBandwidth := []int{1 << 22, 1 << 21, 1 << 21}
	for i := int64(1); i <= n; i++ {
		regionIndex := 0
		r := rand.Float32()
		acc := float32(0)
		for index, f := range nodeDistribution {
			if r > acc && r < acc+f {
				regionIndex = index
			}
			acc += f
		}
		//net.indexes = append(net.indexes, hash.Hash64(uint64(i)))
		net.Add(newNode(
			i,
			downloadBandwidth[regionIndex],
			uploadBandwidth[regionIndex],
			net.Regions[regionIndex],
		), uint64(i))
	}
	//fmt.Println("indexes", net.indexes)
}

//n: at most n connection per node
func (net *Network) InitFloodConnections(n int) {
	for _, node := range net.nodes {
		connectCount := node.RoutingTableLength()
		for connectCount < n {
			r := rand.Intn(net.Size()) + 1
			//fmt.Println(r)
			net.Connect(node, net.nodes[uint64(r)], NewBasicPeerInfo)
			connectCount++
		}
	}
}
func (net *Network) InitKademliaConnections() {
	for _, node := range net.nodes {
		for _, peer := range net.nodes {
			net.Connect(node, peer, NewKadcastPeerInfo)
		}
	}
}

func (net *Network) Node(id uint64) *node.Node {
	return net.nodes[id]
}
func (net *Network) NodeID(id uint64) uint64 {
	return net.indexes[id]
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
	//fmt.Println("connect", a.Id(), b.Id())
	return true
}

func (net *Network) Add(n *node.Node, i uint64) {
	net.nodes[n.Id()] = n
	net.indexes[i] = n.Id()
}

func (net *Network) Size() int {
	return len(net.indexes)
}

//type InfoSorter interface {
//	Append(*information.Information)
//	Take() *information.Information
//}
