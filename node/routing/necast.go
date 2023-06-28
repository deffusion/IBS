package routing

import (
	"math"
	"math/rand"
)

type NeCastTable struct {
	KadcastTable
}

func NewNecastTable(nodeID uint64, k, beta int) Table {
	return &NeCastTable{
		KadcastTable{
			NewKademlia(nodeID, k),
			beta,
			0,
		},
	}
}

func (t *NeCastTable) PeersToBroadcast(from uint64) []uint64 {
	b := 0
	if from != 0 {
		_b, _ := t.kademlia.Locate(from)
		b = _b + 1
	}
	var peers []uint64
	t.SortPeers()
	// broadcast to all peers in buckets of subtree that height less than b
	for i := b; i < KeySpaceBits; i++ {
		ps := t.RandomPeerBasedOnScore(i, t.Beta)
		peers = append(peers, ps...)
	}
	return peers
}

func (t *NeCastTable) IsNeighbour(ID uint64) bool {
	if ID == 0 {
		return false
	}
	_, i := t.Locate(ID)
	if i != -1 {
		return true
	}
	return false
}

func (t *NeCastTable) necastPeerInfo(ID uint64) *NecastPeerInfo {
	b, i := t.Locate(ID)
	if i == -1 {
		return nil
	}
	return t.buckets[b][i].(*NecastPeerInfo)
}

func randomPeersBasedOnScore(peers PeerInfos, n int) []uint64 {
	//if n > len(peers) {
	//	n = len(peers)
	//}
	//totalScore := float64(0)
	//var scores []float64
	//var peerIDS []uint64
	//for _, peer := range peers {
	//	peerIDS = append(peerIDS, peer.PeerID())
	//	scores = append(scores, peer.Score())
	//	totalScore += peer.Score()
	//}
	l := len(peers)
	if n > l {
		n = l
	}
	totalScore := float64(0)
	var scores []float64
	var peerIDS []uint64
	R := int(math.Floor(math.Log2(float64(l))))
	for i, peer := range peers {
		r := int(math.Floor(math.Log2(float64(i + 1))))
		peerIDS = append(peerIDS, peer.PeerID())
		score := 1 << (R - r)
		scores = append(scores, float64(score))
		totalScore += float64(score)
	}
	var randomPeers []uint64
	for n > 0 {
		n--
		nextIndex := 0
		acc := float64(0)
		r := rand.Float64()
		for index, s := range scores {
			if r > acc && r < acc+s/totalScore {
				nextIndex = index
				break
			}
			acc += s / totalScore
		}
		randomPeers = append(randomPeers, peerIDS[nextIndex])
		totalScore -= scores[nextIndex]
		for i := nextIndex; i < len(scores)-1; i++ {
			peerIDS[i] = peerIDS[i+1]
			scores[i] = scores[i+1]
		}
		scores = scores[:len(scores)-1]
		peerIDS = peerIDS[:len(peerIDS)-1]
	}
	return randomPeers
}

func (t *NeCastTable) RandomPeerBasedOnScore(bucket, n int) []uint64 {
	peers := t.buckets[bucket]
	return randomPeersBasedOnScore(peers, n)
}

func (t *NeCastTable) IncrementNewMsg(ID uint64) {
	pi := t.necastPeerInfo(ID)
	if pi == nil {
		return
	}
	pi.NewMsg()
}
func (t *NeCastTable) IncrementConfirmation(ID uint64) {
	pi := t.necastPeerInfo(ID)
	if pi == nil {
		return
	}
	pi.Confirmation()
}
func (t *NeCastTable) IncrementReceivedConfirmation(ID uint64) {
	//fmt.Println("confirm")
	pi := t.necastPeerInfo(ID)
	if pi == nil {
		return
	}
	pi.ReceivedConfirmation()
}
