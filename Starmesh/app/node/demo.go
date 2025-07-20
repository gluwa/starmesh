package node

import (
	"encoding/json"
	"fmt"
	"net"
	"time"

	"starmesh/common/model"
)

// TestConnectionDemo runs a client-server connection test.
func TestConnectionDemo() {
	n := &Node{
		Port:  "8081",
		Peers: make(map[string]model.Peer),
	}

	go n.StartServer()

	time.Sleep(1 * time.Second)

	ConnectToPeer("localhost:8081", "Hello from Node B")
}

// TestRegisterPeerDemo registers sample peers and prints them.
func TestRegisterPeerDemo() {
	n := Node{
		Port:  "8081",
		Peers: make(map[string]model.Peer),
	}

	n.RegisterPeer(model.Peer{
		IP:   "192.168.0.5",
		Port: "8083",
		Lat:  37.5,
		Lon:  127.0,
	})

	n.RegisterPeer(model.Peer{
		IP:   "192.168.0.6",
		Port: "8084",
		Lat:  35.2,
		Lon:  129.0,
	})

	for k, v := range n.Peers {
		fmt.Printf("✅ Peer %s at (%f, %f)\n", k, v.Lat, v.Lon)
	}
}

func TestNHop() {
	// A, B, C 노드 설정
	nodeA := &Node{Port: "9001", Peers: make(map[string]model.Peer)}
	nodeB := Node{Port: "9002", Peers: make(map[string]model.Peer)}
	nodeC := &Node{Port: "9003", Peers: make(map[string]model.Peer)}

	go nodeA.StartServer()
	go nodeB.StartServer()
	go nodeC.StartServer()

	time.Sleep(1 * time.Second)

	// A에서 메시지 시작
	route := []model.Peer{
		{IP: "127.0.0.1", Port: "9002"},
		{IP: "127.0.0.1", Port: "9003"},
	}

	msg := NewMessage(
		route,
		"Hello with UUID",
	)

	data, _ := json.Marshal(msg)
	conn, _ := net.Dial("tcp", "127.0.0.1:9002") // Start from B
	defer conn.Close()
	conn.Write(data)

	time.Sleep(1 * time.Second)
}
