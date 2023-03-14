package output

import (
	"IBS/information"
	"IBS/node"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
)

type Packet struct {
	Id                  int
	Timestamp           int64 `json:"timestamp"` // delay(μs) from the generation(timestamp) of information
	PropagationDelay    int32 `json:"propagationDelay"`
	TransmissionDelay   int32 `json:"transmissionDelay"`
	QueuingDelaySending int32 `json:"queuingDelaySending"`
	//QueuingDelayReceiving int32  `json:"queuingDelayReceiving"`
	From       string `json:"from"`
	To         string `json:"to"`
	Hop        int    `json:"hop"`
	Redundancy bool   `json:"redundancy"`
}

func NewPacket(p *information.BasicPacket) *Packet {
	return &Packet{
		p.ID(),
		p.Timestamp(),
		p.PropagationDelay(),
		p.TransmissionDelay(),
		p.QueuingDelaySending(),
		//p.QueuingDelayReceiving(),
		strconv.FormatUint(p.From().Id(), 16),
		strconv.FormatUint(p.To().Id(), 16),
		p.Hop(),
		p.Redundancy(),
	}
}
func WritePackets(p []*Packet) {
	b, err := json.Marshal(p)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(string(b))
	//os.Create("packets.json")
	err = ioutil.WriteFile("output/output_packets.json", b, 0777)
	if err != nil {
		fmt.Println(err)
	}
}

type Node struct {
	Id                string `json:"id"`
	Region            string `json:"region"`
	DownloadBandwidth int    `json:"downloadBandwidth"` // byte/s
	UploadBandwidth   int    `json:"uploadBandwidth"`
	Running           bool   `json:"running"`
}

func NewNode(n node.Node) *Node {
	return &Node{
		strconv.FormatUint(n.Id(), 16),
		n.Region(),
		n.DownloadBandwidth(),
		n.UploadBandwidth(),
		n.Running(),
	}
}

func WriteNodes(data []*Node) {
	b, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(string(b))
	//os.Create("packets.json")
	err = ioutil.WriteFile("output/output_nodes.json", b, 0777)
	if err != nil {
		fmt.Println(err)
	}
}
