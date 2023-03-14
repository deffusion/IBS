package node

import (
	"IBS/node/routing"
	"fmt"
)

type NeNode struct {
	*BasicNode
	Tasks map[int]*BroadcastTask // metaInfo id -> task
}

func NewNeNode(id uint64, downloadBandwidth, uploadBandwidth int, region string, table routing.Table) *NeNode {
	n := &NeNode{
		NewBasicNode(id, downloadBandwidth, uploadBandwidth, region, table),
		make(map[int]*BroadcastTask),
	}
	return n
}

func (n *NeNode) PeersFromTask(infoID, bucket int) *[]uint64 {
	var peers []uint64
	if bucket >= 0 && bucket < routing.KeySpaceBits {
		return n.peersFromTaskInBucket(infoID, bucket)
	}
	for b := 0; b < routing.KeySpaceBits; b++ {
		peers = append(peers, *n.peersFromTaskInBucket(infoID, b)...)
	}
	return &peers
}
func (n *NeNode) peersFromTaskInBucket(infoID, bucket int) *[]uint64 {
	var peers []uint64
	table := n.routingTable.(*routing.NeCastTable)
	task := n.Tasks[infoID]
	//fmt.Println("n.Tasks", n.Tasks)
	if task.confirmation[bucket] > 0 {
		num := table.MinFanOut
		if num > len(task.candidates[bucket]) {
			num = len(task.candidates[bucket])
		}
		peers = task.candidates[bucket][0:num]
		task.candidates[bucket] = task.candidates[bucket][num:]
	}
	return &peers
}
func (n *NeNode) Confirm(infoID int, from uint64) {
	if from == 0 {
		return
	}
	fmt.Println("confirm", infoID, "from", from)
	table := n.routingTable.(*routing.NeCastTable)
	task := n.Tasks[infoID]
	task.totalConfirm--
	task.confirmation[table.Locate(from)]--
}
func (n *NeNode) IsNeighbour(ID uint64) bool {
	return n.routingTable.(*routing.NeCastTable).IsNeighbour(ID)
}
