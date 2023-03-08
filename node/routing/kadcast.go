package routing

import (
	"fmt"
	"math"
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

func (t *KadcastTable) SetPeerLimit(k int) {
	t.k = k
}

func (t *KadcastTable) PeerLimit() int {
	return math.MaxInt
}

func (t *KadcastTable) AddPeer(peerInfo PeerInfo) error {
	if t.kademlia.AddPeer(peerInfo) {
		t.peerCount++
	}
	// subtract the number of evicted peers
	//t.peerCount -=
	return nil
}

func (t *KadcastTable) RemovePeer(id uint64) {
	// unused
}

func (t *KadcastTable) PeersToBroadcast(from uint64) []uint64 {
	b := 0
	if from != 0 {
		b = t.kademlia.Locate(from) + 1
	}
	//fmt.Println("start from bucket", b)
	//t.PrintTable()
	var peers []uint64
	// broadcast to all peers in buckets of subtree that height less than b
	for i := b; i < KeySpaceBits; i++ {
		for _, info := range t.buckets[i] {
			peers = append(peers, info.PeerID())
		}
	}
	return peers
}

func (t *KadcastTable) SetLastSeen(id uint64, timestamp int64) {
	err := t.kademlia.SetLastSeen(id, timestamp)
	if err != nil {
		fmt.Println(err)
	}
}
func (t *KadcastTable) PrintTable() {
	t.kademlia.PrintBuckets()
}
