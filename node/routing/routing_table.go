package routing

type Table interface {
	Length() int
	SetPeerLimit(int)
	PeerLimit() int
	NoRoomForNewPeer(peerID uint64) bool
	AddPeer(PeerInfo) error
	RemovePeer(uint64)
	PeersToBroadcast(from uint64) []uint64
	SetLastSeen(uint64, int64) error // peerID, timestamp
	PrintTable()
}

// PeerInfo : rank from the peer
//type PeerInfo struct {
//	NodeID  int
//}

type PeerInfo interface {
	PeerID() uint64
	Score() float64 // higher score, higher priority
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

func (i *BasicPeerInfo) Score() float64 {
	return float64(i.LastSeen())
}

func (i *BasicPeerInfo) SetLastSeen(lastSeen int64) {
	i.lastSeen = lastSeen
}
func (i *BasicPeerInfo) LastSeen() int64 {
	return i.lastSeen
}

type NecastPeerInfo struct {
	*BasicPeerInfo
	newMsg               int
	confirmation         int
	receivedConfirmation int
	delay                int32
}

func NewNecastPeerInfo(id uint64) *NecastPeerInfo {
	return &NecastPeerInfo{
		NewBasicPeerInfo(id),
		0,
		0,
		1,
		1,
	}
}
func (n *NecastPeerInfo) SetDelay(delay int32) {
	//n.delay = 1
}
func (n *NecastPeerInfo) NewMsg() {
	n.newMsg += 100
}
func (n *NecastPeerInfo) Confirmation() {
	n.confirmation += 100
}
func (n *NecastPeerInfo) ReceivedConfirmation() {
	n.receivedConfirmation += 100
}
func (n *NecastPeerInfo) Score() float64 {
	return float64(n.newMsg+n.confirmation+n.receivedConfirmation) / float64(n.delay)
}
