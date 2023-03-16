package routing

import (
	"math"
	"sort"
)

type KadcastTable struct {
	*kademlia
	peerCount int
}

func NewKadcastTable(nodeID uint64, k int) *KadcastTable {
	return &KadcastTable{
		NewKademlia(nodeID, k),
		0,
	}
}

func (t *KadcastTable) Length() int {
	return t.peerCount
}

func (t *KadcastTable) NoRoomForNewPeer(peerID uint64) bool {
	b, _ := t.Locate(peerID)
	if len(t.buckets[b]) >= t.k {
		return true
	}
	return false
}

func (t *KadcastTable) SetPeerLimit(k int) {
	t.k = k
}

func (t *KadcastTable) PeerLimit() int {
	return math.MaxInt
}

func (t *KadcastTable) AddPeer(peerInfo PeerInfo) error {
	return t.kademlia.AddPeer(peerInfo)
}

func (t *KadcastTable) RemovePeer(peerInfo PeerInfo) {
	t.kademlia.RemovePeer(peerInfo)
}

func (t *KadcastTable) PeersToBroadcast(from uint64) []uint64 {
	b := 0
	if from != 0 {
		_b, _ := t.kademlia.Locate(from)
		b = _b + 1
	}
	//fmt.Println("start from bucket", b)
	//t.PrintTable()
	var peers []uint64
	// broadcast to all peers in buckets of subtree that height less than b
	for i := b; i < KeySpaceBits; i++ {
		for ind, info := range t.buckets[i] {
			if ind > t.k {
				break
			}
			peers = append(peers, info.PeerID())
		}
	}
	return peers
}

func (t *KadcastTable) SetLastSeen(id uint64, timestamp int64) error {
	return t.kademlia.SetLastSeen(id, timestamp)
}

func PeersInBucket(t *KadcastTable, i int) *[]uint64 {
	return t.kademlia.PeersInBucket(i)
}

func (t *KadcastTable) PrintTable() {
	t.kademlia.PrintBuckets()
}
func (t *KadcastTable) SortPeers() {
	for _, bucket := range t.buckets {
		sort.Sort(bucket)
	}
}
