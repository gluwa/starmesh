package test

import (
	"fmt"
	"starmesh/app/node"
	"starmesh/common/model"
	"testing"
)

func TestRegisterPeer(t *testing.T) {
	n := node.Node{
		Port:  "8081",
		Peers: make(map[string]model.Peer),
	}

	peer1 := model.Peer{
		IP:   "192.168.0.5",
		Port: "8083",
		Lat:  37.5,
		Lon:  127.0,
	}

	peer2 := model.Peer{
		IP:   "192.168.0.6",
		Port: "8084",
		Lat:  35.2,
		Lon:  129.0,
	}

	n.RegisterPeer(peer1)
	n.RegisterPeer(peer2)

	if len(n.Peers) != 2 {
		t.Errorf("expected 2 peers, got %d", len(n.Peers))
	}

	for k, v := range n.Peers {
		fmt.Printf("✅ Peer %s at (%f, %f)\n", k, v.Lat, v.Lon)
	}
}
