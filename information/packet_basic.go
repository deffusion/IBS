package information

import (
	"github.com/deffusion/IBS/node"
)

type BasicPacket struct {
	Information
	timestamp           int64 // delay(μs) from the generation(timestamp) of information
	propagationDelay    int32
	transmissionDelay   int32
	queuingDelaySending int32

	from       node.Node
	to         node.Node
	hop        int
	redundancy bool
}

func NewBasicPacket(id, dataSize int, origin, from, to, relayer node.Node, timestamp int64) *BasicPacket {
	return &BasicPacket{
		Information{
			&meta{
				id,
				timestamp,
				dataSize,
				//net,
				origin,
			},
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

func (p *BasicPacket) NextPacket(to node.Node, propagationDelay, transmissionDelay int32, setRelay bool) *BasicPacket {
	packet := *p
	packet.from = p.to // the last receiver is the next sender
	packet.to = to
	packet.hop++
	packet.propagationDelay = propagationDelay
	packet.transmissionDelay = transmissionDelay
	packet.queuingDelaySending = 0
	packet.timestamp += int64(propagationDelay + transmissionDelay)
	if setRelay {
		packet.relayNode = to
	}
	return &packet
}
func (p *BasicPacket) ConfirmPacket() Packet {
	return NewBasicPacket(p.ID(), 20, p.Origin(), p.To(), p.Origin(), p.Relay(), p.timestamp)
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
func (p *BasicPacket) SetAndAddQueuingDelay(queueing int32) {
	p.queuingDelaySending = queueing
	p.timestamp += int64(queueing)
}

//	func (p *BasicPacket) QueuingDelayReceiving() int32 {
//		return p.queuingDelayReceiving
//	}
func (p *BasicPacket) From() node.Node {
	return p.from
}
func (p *BasicPacket) To() node.Node {
	return p.to
}
func (p *BasicPacket) Redundancy() bool {
	return p.redundancy
}
func (p *BasicPacket) SetRedundancy(r bool) {
	p.redundancy = r
}
func (p *BasicPacket) Hop() int {
	return p.hop
}
