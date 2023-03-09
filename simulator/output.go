package main

import "IBS/information"

type OutputPacket struct {
	Id                  int
	Timestamp           int64 `json:"timestamp"` // delay(μs) from the generation(timestamp) of information
	PropagationDelay    int32 `json:"propagationDelay"`
	TransmissionDelay   int32 `json:"transmissionDelay"`
	QueuingDelaySending int32 `json:"queuingDelaySending"`
	//QueuingDelayReceiving int32  `json:"queuingDelayReceiving"`
	From       uint64 `json:"from"`
	To         uint64 `json:"to"`
	hop        int    `json:"hop"`
	redundancy bool   `json:"redundancy"`
}

func NewOutputPacket(p *information.Packet) *OutputPacket {
	return &OutputPacket{
		p.ID(),
		p.Timestamp(),
		p.PropagationDelay(),
		p.TransmissionDelay(),
		p.QueuingDelaySending(),
		//p.QueuingDelayReceiving(),
		p.From().Id(),
		p.To().Id(),
		p.Hop(),
		p.Redundancy(),
	}
}
