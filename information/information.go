package information

import (
	"IBS/network"
	"IBS/node"
	"fmt"
	"sort"
)

type Information struct {
	id         int
	timestamp  int64
	originNode *node.Node
	dataSize   int // Byte
	net        *network.Network
}

func (i *Information) ID() int {
	return i.id
}

//func (i *Information) getTimestamp() int64 {
//	return i.timestamp
//}
func (i *Information) DataSize() int {
	return i.dataSize
}

type Packet struct {
	*Information
	timestamp           int64 // delay(μs) from the generation(timestamp) of information
	propagationDelay    int32
	transmissionDelay   int32
	queuingDelaySending int32
	//queuingDelayReceiving int32
	from       *node.Node
	to         *node.Node
	hop        int
	redundancy bool
}

type Packets []*Packet

func (ps Packets) Len() int {
	return len(ps)
}
func (ps Packets) Swap(i, j int) {
	ps[i], ps[j] = ps[j], ps[i]
}
func (ps Packets) Less(i, j int) bool {
	return ps[i].timestamp < ps[j].timestamp
}

func (p *Packet) Print() {
	fmt.Printf(
		"pacekt: %d %d->%d originNode: %d size: %dB timestamp: %d propagationDelay: %d transmissionDelay: %d queuingDelaySending: %d\n",
		p.id, p.from.Id(), p.to.Id(), p.originNode.Id(), p.dataSize, p.timestamp, p.propagationDelay, p.transmissionDelay, p.queuingDelaySending)
}

func NewPacket(id, dataSize int, from, to, originNode *node.Node, timestamp int64, net *network.Network) *Packet {
	return &Packet{
		&Information{id, timestamp, originNode, dataSize, net},
		timestamp,
		0,
		0,
		0,
		from,
		to,
		0,
		false,
	}
}

func (p *Packet) NextPacket(to *node.Node, propagationDelay, transmissionDelay int32) *Packet {
	packet := *p
	packet.from = p.to // the last receiver is the next sender
	packet.to = to
	packet.hop++
	packet.propagationDelay = propagationDelay
	packet.transmissionDelay = transmissionDelay
	packet.queuingDelaySending = 0
	packet.timestamp += int64(propagationDelay + transmissionDelay)
	// receiving queue delay
	//if to.TsLastReceived > packet.timestamp {
	//	packet.queuingDelayReceiving = int32(to.TsLastReceived - packet.timestamp)
	//	fmt.Printf("%d->%d(TsLastReceived at %d) queuingDelayReceiving=%d\n",
	//		packet.from.Id(), packet.to.Id(), to.TsLastReceived, packet.queuingDelayReceiving)
	//}
	// sending queuing delay will be considered later
	return &packet
}
func (p *Packet) NextPackets() *Packets {
	var packets Packets
	sender := p.to
	receivedAt := p.timestamp
	received := sender.Received(p.id, p.timestamp)
	if received == true {
		p.redundancy = true
		return &packets
	}

	//fmt.Printf("%d->%d info=%d hop=%d t=%d μs\n", p.from.Id(), sender.Id(), p.id, p.hop, p.timestamp)
	IDs := sender.PeersToBroadCast(p.from)
	regionID := p.net.RegionId
	for _, toID := range *IDs {
		to := p.net.Node(toID)
		// p.to: sender of next packets
		propagationDelay := (*p.net.DelayOfRegions)[regionID[sender.Region()]][regionID[to.Region()]]
		bandwidth := sender.UploadBandwidth()
		if bandwidth > to.DownloadBandwidth() {
			bandwidth = to.DownloadBandwidth()
		}
		transmissionDelay := p.dataSize * 1_000_000 / bandwidth // μs
		packet := p.NextPacket(to, propagationDelay, int32(transmissionDelay))
		packets = append(packets, packet)
	}
	// add sending queuing delay for each packet
	// sending the packet that is earliest to be received first
	sort.Sort(packets)
	base := int32(0)
	if receivedAt < sender.TsLastSending {
		base = int32(sender.TsLastSending - receivedAt)
	}
	for _, packet := range packets {
		packet.queuingDelaySending = base
		packet.timestamp += int64(base)
		base += packet.transmissionDelay
		//packet.to.TsLastReceived = packet.timestamp
	}
	sender.TsLastSending = receivedAt + int64(base)
	return &packets
}

// InfoTimestamp : the time when information in the packet was generated
func (p *Packet) InfoTimestamp() int64 {
	return p.Information.timestamp
}
func (p *Packet) Timestamp() int64 {
	return p.timestamp
}
func (p *Packet) PropagationDelay() int32 {
	return p.propagationDelay
}
func (p *Packet) TransmissionDelay() int32 {
	return p.transmissionDelay
}
func (p *Packet) QueuingDelaySending() int32 {
	return p.queuingDelaySending
}

//func (p *Packet) QueuingDelayReceiving() int32 {
//	return p.queuingDelayReceiving
//}
func (p *Packet) From() *node.Node {
	return p.from
}
func (p *Packet) To() *node.Node {
	return p.to
}
func (p *Packet) Redundancy() bool {
	return p.redundancy
}
func (p *Packet) Hop() int {
	return p.hop
}
