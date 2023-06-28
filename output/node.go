package output

import (
	"IBS/node"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
)

type Node struct {
	Id     string `json:"id"`
	Region string `json:"region"`
	//DownloadBandwidth int    `json:"downloadBandwidth"` // byte/s
	UploadBandwidth int  `json:"uploadBandwidth"`
	Running         bool `json:"running"`
	CrashFactor     int  `json:"crashFactor"`
	CrashTimes      int  `json:"crashTimes"`
}

func newNode(n node.Node) *Node {
	return &Node{
		strconv.FormatUint(n.Id(), 10),
		n.Region(),
		//n.DownloadBandwidth(),
		n.UploadBandwidth(),
		n.Running(),
		n.CrashFactor(),
		n.CrashTimes(),
	}
}

type NodeOutput []*Node

func NewNodeOutput() NodeOutput {
	return NodeOutput{}
}

func (o *NodeOutput) Append(n node.Node) {
	*o = append(*o, newNode(n))
}

func (o *NodeOutput) WriteNodes() {
	b, err := json.Marshal(o)
	if err != nil {
		fmt.Println(err)
	}
	err = ioutil.WriteFile("output/output_nodes.json", b, 0777)
	if err != nil {
		fmt.Println(err)
	}
}
