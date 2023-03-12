package network

import (
	"IBS/node"
	"IBS/node/routing"
	"encoding/json"
	"io/ioutil"
	"math/rand"
)

const BootNodeID = 0

func NewBasicPeerInfo(n *node.Node) routing.PeerInfo {
	return routing.NewBasicPeerInfo(n.Id())
}

type Region struct {
	Region            string  `json:"region"`
	UploadBandwidth   int     `json:"uploadBandwidth"`
	DownloadBandwidth int     `json:"downloadBandwidth"`
	Distribution      float32 `json:"distribution"`
}

type Delays [][]int32

type Network struct {
	bootNode       *node.Node
	nodes          map[uint64]*node.Node
	indexes        map[int]uint64
	DelayOfRegions *Delays

	RegionId          map[string]int
	regions           []string
	nodeDistribution  []float32
	uploadBandwidth   []int
	downloadBandwidth []int
}

func NewNetwork(bootNode *node.Node) *Network {
	// unit: μs (0.000,001s)
	net := &Network{
		bootNode,
		map[uint64]*node.Node{},
		map[int]uint64{},
		&Delays{},

		make(map[string]int),
		[]string{},
		[]float32{},
		[]int{},
		[]int{},
	}
	net.loadConf()
	return net
}

func (net *Network) loadConf() {
	// DelayOfRegions
	delay, err := ioutil.ReadFile("conf/delay.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(delay, net.DelayOfRegions)
	if err != nil {
		panic(err)
	}

	// regions
	var regions []Region
	region, err := ioutil.ReadFile("conf/region.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(region, &regions)
	if err != nil {
		panic(err)
	}
	for i, r := range regions {
		net.RegionId[r.Region] = i
		net.regions = append(net.regions, r.Region)
		net.nodeDistribution = append(net.nodeDistribution, r.Distribution)
		net.uploadBandwidth = append(net.uploadBandwidth, 1<<r.UploadBandwidth)
		net.downloadBandwidth = append(net.downloadBandwidth, 1<<r.DownloadBandwidth)
	}
}

func (net *Network) generateNodes(n int, newNode func(int64, int, int, string) *node.Node) {
	for i := 1; i <= n; i++ {
		regionIndex := 0
		r := rand.Float32()
		acc := float32(0)
		for index, f := range net.nodeDistribution {
			if r > acc && r < acc+f {
				regionIndex = index
			}
			acc += f
		}
		net.Add(newNode(
			int64(i),
			net.downloadBandwidth[regionIndex],
			net.uploadBandwidth[regionIndex],
			net.regions[regionIndex],
		), i)
	}
}

func (net *Network) BootNode() *node.Node {
	return net.bootNode
}

func (net *Network) Node(id uint64) *node.Node {
	return net.nodes[id]
}

// NodeID (index in the network ->nodeID)
func (net *Network) NodeID(id int) uint64 {
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

func (net *Network) Add(n *node.Node, i int) {
	net.nodes[n.Id()] = n
	net.indexes[i] = n.Id()
}

func (net *Network) Size() int {
	return len(net.indexes)
}

func (net *Network) NodeCollapse(n int) {
	for i := 1; i <= n; i++ {
		id := net.NodeID(i)
		net.Node(id).Stop()
	}
}
