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
	timestamp             int64 // delay(μs) from the generation(timestamp) of information
	propagationDelay      int64
	transmissionDelay     int64
	queuingDelaySending   int64
	queuingDelayReceiving int64
	from                  *node.Node
	to                    *node.Node
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
		"pacekt: %d %d->%d originNode: %d size: %dB timestamp: %d propagationDelay: %d transmissionDelay: %d queuingDelaySending: %d queuingDelayReceiving: %d\n",
		p.id, p.from.Id(), p.to.Id(), p.originNode.Id(), p.dataSize, p.timestamp, p.propagationDelay, p.transmissionDelay, p.queuingDelaySending, p.queuingDelayReceiving)
}

func NewPacket(id, dataSize int, from, to, originNode *node.Node, timestamp int64, net *network.Network) *Packet {
	return &Packet{
		&Information{id, timestamp, originNode, dataSize, net},
		timestamp,
		0,
		0,
		0,
		0,
		from,
		to,
	}
}

func (p *Packet) NextPacket(to *node.Node, propagationDelay, transmissionDelay int64) *Packet {
	packet := *p
	packet.from = p.to // the last receiver is the next sender
	packet.to = to
	packet.propagationDelay = propagationDelay
	packet.transmissionDelay = transmissionDelay
	// receiving queue delay
	if to.TsLastReceived > packet.timestamp {
		packet.queuingDelayReceiving = packet.timestamp - to.TsLastReceived
	}
	// sending queuing delay will be considered later
	packet.timestamp = packet.timestamp + propagationDelay + transmissionDelay + packet.queuingDelayReceiving
	//to.TsLastReceived = packet.timestamp
	return &packet
}
func (p *Packet) NextPackets() *Packets {
	var packets Packets
	sender := p.to
	receivedAt := p.timestamp
	received := sender.Received(p.id, p.timestamp)
	if received == true {
		return &packets
	} else {
		//fmt.Printf("%d->%d info=%d t=%d μs\n", p.from.Id(), sender.Id(), p.id, p.timestamp)
	}
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
		packet := p.NextPacket(to, propagationDelay, int64(transmissionDelay))
		packets = append(packets, packet)
	}
	// add sending queuing delay for each packet
	// sending the packet that is earliest to be received first
	sort.Sort(packets)
	base := int64(0)
	for i, packet := range packets {
		if i == 0 {
			if sender.TsLastSend > receivedAt {
				base = sender.TsLastSend - receivedAt
			}
		}
		packet.queuingDelaySending = base
		packet.timestamp += (packet.queuingDelaySending - packet.queuingDelayReceiving)
		base += packet.transmissionDelay
	}

	return &packets
}

// InfoTimestamp : the time when information in the packet was generated
func (p *Packet) InfoTimestamp() int64 {
	return p.Information.timestamp
}
func (p *Packet) Timestamp() int64 {
	return p.timestamp
}
func (p *Packet) From() *node.Node {
	return p.from
}
func (p *Packet) To() *node.Node {
	return p.to
}
