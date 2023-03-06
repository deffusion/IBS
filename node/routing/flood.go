package routing

import (
	"errors"
	"fmt"
)

type FloodPeerInfo struct {
	id int
}

func (i *FloodPeerInfo) PeerID() int {
	return i.id
}

func (i *FloodPeerInfo) Score() int {
	return 1
}

func NewFloodPeerInfo(id int) *FloodPeerInfo {
	return &FloodPeerInfo{id}
}

type FloodTable struct {
	// id:int
	table map[int]PeerInfo
	limit int
}

func NewFloodTable(limit int) *FloodTable {
	return &FloodTable{
		make(map[int]PeerInfo),
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

func (t *FloodTable) RemovePeer(id int) {
	delete(t.table, id)
}

func (t *FloodTable) PeersToBroadcast(from int) []int {
	var peers []int
	// broadcast to all peers except the sender
	for id, _ := range t.table {
		if id != from {
			peers = append(peers, id)
		}
	}
	return peers
}
