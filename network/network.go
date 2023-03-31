package network

import (
	"IBS/node"
	"IBS/node/routing"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"time"
)

const BootNodeID = 0

type NewPeerInfo func(node.Node) routing.PeerInfo

func NewBasicPeerInfo(n node.Node) routing.PeerInfo {
	return routing.NewBasicPeerInfo(n.Id())
}

type Region struct {
	Region       string  `json:"region"`
	Distribution float32 `json:"distribution"`
}

type Bandwidth struct {
	UploadBandwidth int     `json:"uploadBandwidth"`
	Distribution    float32 `json:"distribution"`
}

type Delays [][]int32

type Network struct {
	bootNode       node.Node
	Nodes          map[uint64]node.Node
	indexes        map[int]uint64
	DelayOfRegions *Delays

	RegionId                    map[string]int
	regions                     []string
	nodeDistribution            []float32
	uploadBandwidths            []int
	uploadBandwidthDistribution []float32
	//downloadBandwidth []int
}

func NewNetwork(bootNode node.Node) *Network {
	// unit: μs (0.000,001s)
	net := &Network{
		bootNode,
		map[uint64]node.Node{},
		map[int]uint64{},
		&Delays{},

		make(map[string]int),
		[]string{},
		[]float32{},
		[]int{},
		[]float32{},
		//[]int{},
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
		//net.uploadBandwidth = append(net.uploadBandwidth, 1<<r.UploadBandwidth)
		//net.downloadBandwidth = append(net.downloadBandwidth, 1<<r.DownloadBandwidth)
	}

	// bandwidth
	var bandwidths []Bandwidth
	bandwidth, err := ioutil.ReadFile("conf/bandwidth.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(bandwidth, &bandwidths)
	if err != nil {
		panic(err)
	}
	for _, b := range bandwidths {
		net.uploadBandwidthDistribution = append(net.uploadBandwidthDistribution, b.Distribution)
		net.uploadBandwidths = append(net.uploadBandwidths, 1<<b.UploadBandwidth)
	}

	fmt.Println("upload bandwidth:", net.uploadBandwidths)
}

func (net *Network) generateNodes(n int, newNode func(int, int, string, int) node.Node, degree int) {
	rand.Seed(time.Now().Unix())
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
		bandwidthIndex := 0
		r = rand.Float32()
		acc = float32(0)
		for index, f := range net.uploadBandwidthDistribution {
			if r > acc && r < acc+f {
				bandwidthIndex = index
			}
			acc += f
		}
		_node := newNode(
			i,
			//net.downloadBandwidth[regionIndex],
			net.uploadBandwidths[bandwidthIndex],
			net.regions[regionIndex],
			degree,
		)
		net.Add(_node, i)
	}
}

func (net *Network) BootNode() node.Node {
	return net.bootNode
}

func (net *Network) Node(id uint64) node.Node {
	return net.Nodes[id]
}

// NodeID (index in the network ->nodeID)
func (net *Network) NodeID(id int) uint64 {
	return net.indexes[id]
}

//func (net *Network) Connect(a, b node.Node, f func(node.Node) routing.PeerInfo) bool {
//	if a.Id() == b.Id() {
//		return false
//	}
//	//fmt.Printf("connect %d to %d\n", a.Id(), b.Id())
//	bInfo := f(b)
//	if a.AddPeer(bInfo) == false {
//		return false
//	}
//	aInfo := f(a)
//	if b.AddPeer(aInfo) == false {
//		a.RemovePeer(bInfo)
//		return false
//	}
//
//	switch aInfo.(type) {
//	case *routing.NecastPeerInfo:
//		delay := (*net.DelayOfRegions)[net.RegionId[a.Region()]][net.RegionId[b.Region()]]
//		aInfo.(*routing.NecastPeerInfo).SetDelay(delay)
//		bInfo.(*routing.NecastPeerInfo).SetDelay(delay)
//	}
//
//	//fmt.Println("connect", a.Id(), b.Id())
//	return true
//}

func (net *Network) Connect(a, b node.Node, f NewPeerInfo) bool {
	//fmt.Println("connect", a.Id(), b.Id())
	if a.Id() == b.Id() {
		return false
	}
	//if a.NoRoomForNewPeer(b.Id()) || b.NoRoomForNewPeer(a.Id()) {
	//	return false
	//}
	//fmt.Printf("connect %d to %d\n", a.Id(), b.Id())
	bInfo := f(b)
	aInfo := f(a)
	switch bInfo.(type) {
	case *routing.NecastPeerInfo:
		delay := (*net.DelayOfRegions)[net.RegionId[a.Region()]][net.RegionId[b.Region()]]
		aInfo.(*routing.NecastPeerInfo).SetDelay(delay)
		bInfo.(*routing.NecastPeerInfo).SetDelay(delay)
	}
	a.AddPeer(bInfo)
	b.AddPeer(aInfo)
	//fmt.Println("connect", a.Id(), b.Id())
	return true
}

func (net *Network) Add(n node.Node, i int) {
	net.Nodes[n.Id()] = n
	net.indexes[i] = n.Id()
}

func (net *Network) Size() int {
	return len(net.indexes)
}

// NodeCrash crash nodes from i to netSize (according to correspond nodes)
func (net *Network) NodeCrash(i int) int {
	rand.Seed(time.Now().UnixMilli())
	cnt := 0
	if i < 1 {
		i = 1
	}
	for ; i <= net.Size(); i++ {
		id := net.NodeID(i)
		n := net.Node(id)
		r := rand.Intn(net.Size())
		if n.CrashFactor() >= r {
			cnt++
			n.Stop()
		}
	}
	return cnt
}

func (net *Network) NodeInfest(i int) int {
	rand.Seed(time.Now().UnixMilli())
	cnt := 0
	if i < 1 {
		i = 1
	}
	for ; i <= net.Size(); i++ {
		id := net.NodeID(i)
		n := net.Node(id)
		r := rand.Intn(net.Size())
		if n.CrashFactor() >= r {
			cnt++
			n.Infest()
		}
	}
	return cnt
}
