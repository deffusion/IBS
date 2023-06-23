package output

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type info struct {
	Id     int    `json:"id"`
	Delay  int    `json:"delay"`
	Region string `json:"region"`
}

type DelayOutput []*info // packetID -> received

func NewDelayOutput() DelayOutput {
	return DelayOutput{}
}
func (o *DelayOutput) Append(id, delay int, region string) {
	*o = append(*o, &info{id, delay, region})
}

func (o *DelayOutput) WriteDelay() {
	b, err := json.Marshal(o)
	if err != nil {
		fmt.Println(err)
	}
	err = ioutil.WriteFile("output/output_delay.json", b, 0777)
	if err != nil {
		fmt.Println(err)
	}
}
