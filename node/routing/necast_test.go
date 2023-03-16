package routing

import (
	"fmt"
	"testing"
)

func TestProbability(t *testing.T) {
	var peers PeerInfos
	p1 := NewNecastPeerInfo(0)
	p1.SetDelay(100)
	p1.Confirmation()
	p1.ReceivedConfirmation()
	p1.NewMsg()
	p1.Confirmation()
	p1.ReceivedConfirmation()
	p1.NewMsg()
	p2 := NewNecastPeerInfo(1)
	p2.SetDelay(100)
	p2.Confirmation()
	p2.ReceivedConfirmation()
	p2.NewMsg()
	p2.NewMsg()
	p3 := NewNecastPeerInfo(2)
	p3.SetDelay(2)
	p3.Confirmation()
	p3.ReceivedConfirmation()
	p3.NewMsg()
	p4 := NewNecastPeerInfo(3)
	p4.Confirmation()
	p4.ReceivedConfirmation()
	p4.NewMsg()
	p4.ReceivedConfirmation()
	p4.NewMsg()
	p4.SetDelay(1)
	peers = append(peers, p1, p2, p3, p4)
	fmt.Println("before", peers)
	fmt.Println(randomPeersBasedOnScore(peers, 3))
	fmt.Println("after", peers)
}
