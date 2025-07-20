package node

import (
	"encoding/json"
	"fmt"
	"net"

	"starmesh/common/model"
	"starmesh/common/util"
)

type GossipMessage struct {
	Type     string       `json:"type"`      // "push" or "pull"
	SenderID string       `json:"sender_id"` // 보낸 노드의 IP:Port
	Peers    []model.Peer `json:"peers"`     // 전파하고 싶은 Peer 목록 일부
}

func (n *Node) SendGossipTo(peer model.Peer) {
	n.PeersMux.RLock()
	defer n.PeersMux.RUnlock()

	// 본인 제외한 Peer들 중 일부만 선택
	peersToSend := make([]model.Peer, 0, 5)
	for _, p := range n.Peers {
		if p.IP != peer.IP || p.Port != peer.Port { // 자기 자신 제외
			peersToSend = append(peersToSend, p)
			if len(peersToSend) >= 5 {
				break
			}
		}
	}

	msg := GossipMessage{
		Type:     "push",
		SenderID: fmt.Sprintf("%s:%s", util.GetOutboundIP(), n.Port),
		Peers:    peersToSend,
	}

	data, _ := json.Marshal(msg)

	address := peer.IP + ":" + peer.Port
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println("[Gossip] Failed to send to", address)
		return
	}
	defer conn.Close()

	conn.Write(data)
	fmt.Printf("[Gossip] Sent to %s (%d peers)\n", address, len(peersToSend))
}
