package output

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type DelayOutput map[int]int // packetID -> received

func NewDelayOutput() *DelayOutput {
	c := make(DelayOutput)
	return &c
}
func (o *DelayOutput) WriteDelay() {
	var outputs [2][]int
	for id, delay := range *o {
		outputs[0] = append(outputs[0], id)
		outputs[1] = append(outputs[1], delay)
	}
	b, err := json.Marshal(outputs)
	if err != nil {
		fmt.Println(err)
	}
	err = ioutil.WriteFile("output/output_delay.json", b, 0777)
	if err != nil {
		fmt.Println(err)
	}
}
