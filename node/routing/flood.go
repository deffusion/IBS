package routing

import (
	"errors"
	"fmt"
)

type FloodTable struct {
	// id:int
	table map[uint64]PeerInfo
	limit int
}

func NewFloodTable(limit int) *FloodTable {
	return &FloodTable{
		make(map[uint64]PeerInfo),
		limit,
	}
}

func (t *FloodTable) Length() int {
	return len(t.table)
}

func (t *FloodTable) SetPeerLimit(n int) {
	t.limit = n
}

func (t *FloodTable) PeerLimit() int {
	return t.limit
}

func (t *FloodTable) AddPeer(peerInfo PeerInfo) error {
	if len(t.table) < t.limit {
		t.table[peerInfo.PeerID()] = peerInfo
	} else {
		s := fmt.Sprintf("adding peer into a full table, size:%d", t.limit)
		return errors.New(s)
	}
	return nil
}

func (t *FloodTable) RemovePeer(id uint64) {
	delete(t.table, id)
}

func (t *FloodTable) PeersToBroadcast(from uint64) []uint64 {
	var peers []uint64
	// broadcast to all peers except the sender
	for id, _ := range t.table {
		if uint64(id) != from {
			peers = append(peers, uint64(id))
		}
	}
	return peers
}

func (t *FloodTable) SetLastSeen(id uint64, timestamp int64) {
	t.table[id].SetLastSeen(timestamp)
}
func (t *FloodTable) PrintTable() {

}
