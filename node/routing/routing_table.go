package routing

type Table interface {
	Length() int
	SetPeerLimit(int)
	PeerLimit() int
	AddPeer(PeerInfo) error
	RemovePeer(uint64)
	PeersToBroadcast(from uint64) []uint64
	SetLastSeen(uint64, int64) // peerID, timestamp
	PrintTable()
}

// PeerInfo : rank from the peer
//type PeerInfo struct {
//	NodeID  int
//}

type PeerInfo interface {
	PeerID() uint64
	Score() int64 // higher score, higher priority
	SetLastSeen(int64)
	LastSeen() int64
}

type PeerInfos []PeerInfo

func (ps PeerInfos) Len() int {
	return len(ps)
}
func (ps PeerInfos) Swap(i, j int) {
	ps[i], ps[j] = ps[j], ps[i]
}
func (ps PeerInfos) Less(i, j int) bool {
	return ps[i].Score() > ps[j].Score()
}

func (ps PeerInfos) Includes(p PeerInfo) bool {
	for i := 0; i < ps.Len(); i++ {
		if ps[i].PeerID() == p.PeerID() {
			return true
		}
	}
	return false
}

type BasicPeerInfo struct {
	id       uint64
	lastSeen int64
}

func NewBasicPeerInfo(id uint64) *BasicPeerInfo {
	return &BasicPeerInfo{id, 0}
}

func (i *BasicPeerInfo) PeerID() uint64 {
	return i.id
}

func (i *BasicPeerInfo) Score() int64 {
	return i.LastSeen()
}

func (i *BasicPeerInfo) SetLastSeen(lastSeen int64) {
	i.lastSeen = lastSeen
}
func (i *BasicPeerInfo) LastSeen() int64 {
	return i.lastSeen
}
