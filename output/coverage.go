package output

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type PacketCoverageOutput map[int]int // packetID -> received

func NewCoverageOutput() PacketCoverageOutput {
	return make(PacketCoverageOutput)
}
func (o *PacketCoverageOutput) WriteCoverage() {
	var outputs [2][]int
	for id, received := range *o {
		outputs[0] = append(outputs[0], id)
		outputs[1] = append(outputs[1], received)
	}
	b, err := json.Marshal(outputs)
	if err != nil {
		fmt.Println(err)
	}
	err = ioutil.WriteFile("output/output_coverage.json", b, 0777)
	if err != nil {
		fmt.Println(err)
	}
}
