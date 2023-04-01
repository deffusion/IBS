package information

import (
	"IBS/network"
	"IBS/node"
	"sort"
)

type BasicPacket struct {
	Information
	timestamp           int64 // delay(μs) from the generation(timestamp) of information
	propagationDelay    int32
	transmissionDelay   int32
	queuingDelaySending int32
	//queuingDelayReceiving int32
	from       node.Node
	to         node.Node
	hop        int
	redundancy bool
}

func NewBasicPacket(id, dataSize int, origin, from, to, relayer node.Node, timestamp int64, net *network.Network) *BasicPacket {
	return &BasicPacket{
		Information{
			&meta{
				id,
				timestamp,
				dataSize,
				net,
			},
			origin,
			relayer,
		},
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

func (p *BasicPacket) nextPacket(to node.Node, propagationDelay, transmissionDelay int32) Packet {
	packet := *p
	packet.from = p.to // the last receiver is the next sender
	packet.to = to
	packet.hop++
	packet.propagationDelay = propagationDelay
	packet.transmissionDelay = transmissionDelay
	packet.queuingDelaySending = 0
	packet.timestamp += int64(propagationDelay + transmissionDelay)
	return &packet
}
func (p *BasicPacket) ConfirmPacket() Packet {
	return NewBasicPacket(p.ID(), 20, p.Origin(), p.To(), p.Origin(), p.Relay(), p.timestamp, p.net)
}
func (p *BasicPacket) NextPackets(IDs *[]uint64) Packets {
	var packets Packets
	sender := p.to
	if sender.Running() == false {
		return packets
	}
	if sender.Malicious() == true {
		p.redundancy = true
		return packets
	}
	receivedAt := p.timestamp
	received := sender.Received(p.id, p.timestamp)
	if received == true {
		p.redundancy = true
		//fmt.Printf("%d->%d info=%d hop=%d t=%d μs (redundancy: %t)\n", p.from.Id(), sender.Id(), p.id, p.hop, p.timestamp, p.redundancy)
		return packets
	}
	switch sender.(type) {
	case *node.NeNode:
		sender.(*node.NeNode).NewMsg(p.From().Id())
	}
	//fmt.Printf("%d->%d info=%d hop=%d t=%d μs (redundancy: %t)\n", p.from.Id(), sender.Id(), p.id, p.hop, p.timestamp, p.redundancy)
	//IDs := sender.PeersToBroadCast(p.from)
	regionID := p.net.RegionId
	for _, toID := range *IDs {
		to := p.net.Node(toID)
		if to.Running() == false {
			continue
		}
		// p.to: sender of next packets
		propagationDelay := (*p.net.DelayOfRegions)[regionID[sender.Region()]][regionID[to.Region()]]
		bandwidth := sender.UploadBandwidth()
		transmissionDelay := p.dataSize * 1_000_000 / bandwidth // μs
		packet := p.nextPacket(to, propagationDelay, int32(transmissionDelay))
		//if p.from.Id() == p.net.BootNode().Id() {
		//	packet.relayer = to
		//}
		//log.Println("fromID:", p.From().Id())
		if p.From().Id() == network.BootNodeID {
			//log.Println("set relayNode", to.Id())
			packet.(*BasicPacket).relayNode = to
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
		packet.(*BasicPacket).queuingDelaySending = base
		packet.(*BasicPacket).timestamp += int64(base)
		base += packet.(*BasicPacket).transmissionDelay
		//packet.to.TsLastReceived = packet.timestamp
	}
	sender.SetTsLastSending(receivedAt + int64(base))
	return packets
}

// InfoTimestamp : the time when information in the packet was generated
func (p *BasicPacket) InfoTimestamp() int64 {
	return p.Information.timestamp
}
func (p *BasicPacket) Timestamp() int64 {
	return p.timestamp
}
func (p *BasicPacket) Origin() node.Node {
	return p.originNode
}
func (p *BasicPacket) Relay() node.Node {
	return p.relayNode
}
func (p *BasicPacket) PropagationDelay() int32 {
	return p.propagationDelay
}
func (p *BasicPacket) TransmissionDelay() int32 {
	return p.transmissionDelay
}
func (p *BasicPacket) QueuingDelaySending() int32 {
	return p.queuingDelaySending
}

//func (p *BasicPacket) QueuingDelayReceiving() int32 {
//	return p.queuingDelayReceiving
//}
func (p *BasicPacket) From() node.Node {
	return p.from
}
func (p *BasicPacket) To() node.Node {
	return p.to
}
func (p *BasicPacket) Redundancy() bool {
	return p.redundancy
}
func (p *BasicPacket) Hop() int {
	return p.hop
}
