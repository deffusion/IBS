package network

import (
	"IBS/information"
	"IBS/node"
	"IBS/node/routing"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"sort"
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
	indexes        map[int]uint64 // order in the network -> id
	DelayOfRegions Delays

	RegionId                    map[string]int
	regions                     []string
	nodeDistribution            []float32
	uploadBandwidths            []int
	uploadBandwidthDistribution []float32
}

func NewNetwork(bootNode node.Node) *Network {
	// unit: μs (0.000,001s)
	net := &Network{
		bootNode,
		map[uint64]node.Node{},
		map[int]uint64{},
		Delays{},

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
	err = json.Unmarshal(delay, &net.DelayOfRegions)
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

// generateNodes generate nodes by given newNode function, its region and bandwidth
// is randomly assigned according to configuration files. And add the node into network
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

// Connect two peers
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

func (net *Network) succeedingPackets(p *information.BasicPacket, IDs *[]uint64) information.Packets {
	var packets information.Packets
	sender := p.To()
	if sender.Running() == false {
		return packets
	}
	if sender.Malicious() == true {
		p.SetRedundancy(true)
		return packets
	}
	receivedAt := p.Timestamp()
	received := sender.Received(p.ID(), p.Timestamp())
	if received == true {
		p.SetRedundancy(true)
		//fmt.Printf("%d->%d info=%d hop=%d t=%d μs (redundancy: %t)\n", p.from.Id(), sender.Id(), p.id, p.hop, p.timestamp, p.redundancy)
		return packets
	}
	switch sender.(type) {
	case *node.NeNode:
		sender.(*node.NeNode).NewMsg(p.From().Id())
	}
	//fmt.Printf("%d->%d info=%d hop=%d t=%d μs (redundancy: %t)\n", p.from.Id(), sender.Id(), p.id, p.hop, p.timestamp, p.redundancy)
	//IDs := sender.PeersToBroadCast(p.from)
	regionID := net.RegionId
	for _, toID := range *IDs {
		to := net.Node(toID)
		if to.Running() == false {
			continue
		}
		// p.to: sender of next packets
		propagationDelay := net.DelayOfRegions[regionID[sender.Region()]][regionID[to.Region()]]
		bandwidth := sender.UploadBandwidth()
		transmissionDelay := p.DataSize() * 1_000_000 / bandwidth // μs
		var packet *information.BasicPacket

		//if p.from.Id() == p.net.BootNode().Id() {
		//	packet.relayer = to
		//}
		//log.Println("fromID:", p.From().Id())
		if p.From().Id() == BootNodeID {
			//log.Println("set relayNode", to.Id())
			packet = p.NextPacket(to, propagationDelay, int32(transmissionDelay), true)
		} else {
			packet = p.NextPacket(to, propagationDelay, int32(transmissionDelay), false)
		}
		packets = append(packets, packet)
	}
	// add sending queuing delay for each packet
	// sending the packet that is earliest to be received first
	sort.Sort(packets)
	base := int32(0)
	if receivedAt < sender.TsLastSending() {
		base = int32(sender.TsLastSending() - receivedAt)
	}
	for _, packet := range packets {
		bp := packet.(*information.BasicPacket)
		bp.SetAndAddQueuingDelay(base)
		base += bp.TransmissionDelay()
		//packet.to.TsLastReceived = packet.timestamp
	}
	sender.SetTsLastSending(receivedAt + int64(base))
	return packets
}

func (net *Network) PacketReplacement(p *information.BasicPacket, confirmCnt *int) information.Packets {
	var packets information.Packets
	var peers = p.To().PeersToBroadCast(p.From())
	switch p.To().(type) {
	case *node.NeNode:
		neNode := p.To().(*node.NeNode)
		//peers = neNode.PeersToBroadCast(p.From())
		if neNode.Id() != p.Origin().Id() && neNode.IsNeighbour(p.Origin().Id()) && neNode.Id() != p.Relay().Id() {
			*confirmCnt++
			// send redundancy confirm packets to the origin node
			packets = append(packets, p.ConfirmPacket())
		}
	}
	crashCnt := 0
	for i, peerID := range peers {
		peers[i-crashCnt] = peers[i]
		//fmt.Println("peerID", peerID)
		if net.Node(peerID).Running() == false {
			p.To().RemovePeer(peerID)
			crashCnt++
		}
	}
	peers = peers[:len(peers)-crashCnt]

	return append(net.succeedingPackets(p, &peers), packets...)
}
