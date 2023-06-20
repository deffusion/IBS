package node

//unused
//type BroadcastTask struct {
//	infoID       int
//	confirmation [routing.KeySpaceBits]int
//	totalConfirm int
//	candidates   [routing.KeySpaceBits][]uint64
//}

//func (n *NeNode) NewBroadcastTask(infoID int) {
//	//n.SortPeers()
//	//n.PrintTable()
//	t := &BroadcastTask{}
//	t.infoID = infoID
//	for b := 0; b < routing.KeySpaceBits; b++ {
//		nTable := n.routingTable.(*routing.NeCastTable)
//		t.candidates[b] = *routing.PeersInBucket(&nTable.KadcastTable, b)
//		t.confirmation[b] = len(t.candidates[b])
//		t.totalConfirm += t.confirmation[b]
//	}
//	n.Tasks[infoID] = t
//}
