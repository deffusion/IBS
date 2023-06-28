package information

//type TimerPacket struct {
//	*meta
//	timestamp int64
//	from      node.Node
//	to        node.Node
//}
//
//func NewTimerPacket(id, dataSize int, from, to node.Node, timestamp int64, net *network.BaseNetwork) *TimerPacket {
//	return &TimerPacket{
//		&meta{
//			id,
//			timestamp,
//			dataSize,
//			net,
//		},
//		timestamp,
//		from,
//		to,
//	}
//}
//
//func (p *TimerPacket) ID() int {
//	return p.id
//}
//
//func (p *TimerPacket) NextPacket(delay int64) Packet {
//	packet := *p
//	packet.timestamp += delay
//	return &packet
//}
//
//func (p *TimerPacket) Timestamp() int64 {
//	return p.timestamp
//}
//
////func (p *BasicPacket) QueuingDelayReceiving() int32 {
////	return p.queuingDelayReceiving
////}
//func (p *TimerPacket) From() node.Node {
//	return p.from
//}
//func (p *TimerPacket) To() node.Node {
//	return p.to
//}
