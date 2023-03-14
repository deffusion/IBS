package routing

type NeCastTable struct {
	MinFanOut int
	KadcastTable
}

func NewNecastTable(nodeID uint64, bucketSize, minFanOut int) *NeCastTable {
	return &NeCastTable{
		minFanOut,
		KadcastTable{
			NewKademlia(nodeID, bucketSize),
			0,
		},
	}
}

func (t *NeCastTable) PeersToBroadcast(from uint64) []uint64 {
	b := 0
	if from != 0 {
		b = t.kademlia.Locate(from) + 1
	}
	var peers []uint64
	// broadcast to all peers in buckets of subtree that height less than b
	for i := b; i < KeySpaceBits; i++ {
		for j, info := range t.buckets[i] {
			if j >= t.MinFanOut {
				break
			}
			peers = append(peers, info.PeerID())
		}
	}
	return peers
}

func (t *NeCastTable) IsNeighbour(ID uint64) bool {
	if ID == 0 {
		return false
	}
	b := t.Locate(ID)
	for _, peerInfo := range t.kademlia.buckets[b] {
		if peerInfo.PeerID() == ID {
			return true
		}
	}
	return false
}
