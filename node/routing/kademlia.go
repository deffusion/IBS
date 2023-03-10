package routing

import (
	"errors"
	"fmt"
	"math/rand"
	"sort"
	"time"
)

const KeySpaceBits = 4

type kademlia struct {
	nodeID  uint64
	buckets [KeySpaceBits]PeerInfos
	k       int // bucket size
}

func NewKademlia(nodeID uint64, k int) *kademlia {
	return &kademlia{
		nodeID,
		[KeySpaceBits]PeerInfos{},
		k,
	}
}

func fakeIDForBucket(num uint64, b int) uint64 {
	var base uint64 = 1 << b
	rand.Seed(time.Now().Unix())
	if b == 63 {
		r := uint64(rand.Int63n(62))
		r += uint64(rand.Intn(2)) << 62 // rand 63bit uint
		return base + r
	}
	return num ^ (uint64(rand.Int63n(int64(base))) + base)
}

func FakeIDForBucket(nodeID uint64, b int) (uint64, error) {
	if b < 0 || b > 63 {
		return nodeID, errors.New("FakeIDForBucket: out of range")
	}
	if b == 0 {
		return nodeID ^ 1, nil
	}
	return fakeIDForBucket(nodeID, b), nil
}

func (k *kademlia) SetLastSeen(id uint64, timestamp int64) error {
	b := k.nodeID ^ id
	for _, info := range k.buckets[b] {
		if info.PeerID() == id {
			info.SetLastSeen(timestamp)
			return nil
		}
	}
	str := fmt.Sprint("kademlia SetLastSeen: failed to find peer", id)
	return errors.New(str)
}

func locate(k1, k2 uint64) int {
	xor := k1 ^ k2
	b := KeySpaceBits
	for xor != 0 {
		xor = xor >> 1
		b--
	}
	return b
}

func (k *kademlia) Locate(peerID uint64) int {
	return locate(k.nodeID, peerID)
}

func (k *kademlia) AddPeer(info PeerInfo) bool {
	b := k.Locate(info.PeerID())
	if k.buckets[b].Includes(info) {
		// already included
		return false
	}
	if k.buckets[b].Len() < k.k {
		k.buckets[b] = append(k.buckets[b], info)
		sort.Sort(k.buckets[b])
		return true
	}
	last := k.buckets[b][k.buckets[b].Len()-1]
	if last.Score() < info.Score() {
		k.buckets[b][k.buckets[b].Len()-1] = info
		sort.Sort(k.buckets[b])
		return true
	}
	return false
}

func (k *kademlia) PrintBuckets() {
	for i, bucket := range k.buckets {
		fmt.Print("bucket", i, ": ")
		for _, info := range bucket {
			fmt.Printf("distance=%d(%d), ", info.PeerID()^k.nodeID, info.PeerID())
		}
		fmt.Println()
	}
}
