package routing

type Table interface {
	Length() int
	SetPeerLimit(int)
	PeerLimit() int
	AddPeer(PeerInfo) error
	RemovePeer(int)
	PeersToBroadcast(from int) []int
}

// PeerInfo : rank from the peer
//type PeerInfo struct {
//	NodeID  int
//}

type PeerInfo interface {
	PeerID() int
	Score() int // higher score, higher priority
}
