package node

import (
	"IBS/node/routing"
)

type NeNode struct {
	*BasicNode
	//Tasks map[int]*BroadcastTask // metaInfo id -> task
}

func NewNeNode(id uint64, downloadBandwidth, uploadBandwidth, crashFactor int, region string, table routing.Table) *NeNode {
	n := &NeNode{
		NewBasicNode(id, downloadBandwidth, uploadBandwidth, crashFactor, region, table),
		//make(map[int]*BroadcastTask),
	}
	return n
}

//func (n *NeNode) PeersFromTask(infoID, bucket int) []uint64 {
//	var peers []uint64
//	if bucket >= 0 && bucket < routing.KeySpaceBits {
//		return n.peersFromTaskInBucket(infoID, bucket)
//	}
//	for b := 0; b < routing.KeySpaceBits; b++ {
//		peers = append(peers, n.peersFromTaskInBucket(infoID, b)...)
//	}
//	return peers
//}
//func (n *NeNode) peersFromTaskInBucket(infoID, bucket int) []uint64 {
//	var peers []uint64
//	table := n.routingTable.(*routing.NeCastTable)
//	task := n.Tasks[infoID]
//	//fmt.Println("n.Tasks", n.Tasks)
//	if task.confirmation[bucket] > 0 {
//		num := table.MinFanOut
//		if num > len(task.candidates[bucket]) {
//			num = len(task.candidates[bucket])
//		}
//		peers = table.RandomPeerBasedOnScore(bucket, num)
//		//TODO: candidate subsection
//		//task.candidates[bucket] = task.candidates[bucket][num:]
//	}
//	return peers
//}
//func (n *NeNode) Confirm(infoID int, from, relay uint64) {
//	if from == 0 || from == relay {
//		return
//	}
//	//fmt.Println("initiator", n.Id(), "confirm", infoID, "from", from, "relay", relay)
//	table := n.routingTable.(*routing.NeCastTable)
//	task := n.Tasks[infoID]
//	task.totalConfirm--
//	b, _ := table.Locate(from)
//	task.confirmation[b]--
//	n.Confirmation(from, relay)
//}
func (n *NeNode) IsNeighbour(ID uint64) bool {
	return n.routingTable.(*routing.NeCastTable).IsNeighbour(ID)
}
func (n *NeNode) NewMsg(peerID uint64) {
	table := n.routingTable.(*routing.NeCastTable)
	table.IncrementNewMsg(peerID)
}
func (n *NeNode) Confirmation(peerID, relayID uint64) {
	table := n.routingTable.(*routing.NeCastTable)
	table.IncrementConfirmation(peerID)
	table.IncrementReceivedConfirmation(peerID)
}

//func (n *NeNode) ReceivedConfirmation(peerID uint64) {
//	table := n.routingTable.(*routing.NeCastTable)
//	table.IncrementReceivedConfirmation(peerID)
//}
func (n *NeNode) SortPeers() {
	table := n.routingTable.(*routing.NeCastTable)
	table.SortPeers()
}
