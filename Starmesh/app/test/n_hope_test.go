package test

import (
	"encoding/json"
	"net"
	"time"

	"starmesh/app/node"
	"starmesh/common/model"
	"testing"
)

func NHop_Test(t *testing.T) {
	// A, B, C 노드 설정
	nodeA := &node.Node{Port: "9001", Peers: make(map[string]model.Peer)}
	nodeB := &node.Node{Port: "9002", Peers: make(map[string]model.Peer)}
	nodeC := &node.Node{Port: "9003", Peers: make(map[string]model.Peer)}

	go nodeA.StartServer()
	go nodeB.StartServer()
	go nodeC.StartServer()

	time.Sleep(1 * time.Second)

	// A에서 메시지 시작
	route := []model.Peer{
		{IP: "127.0.0.1", Port: "9002"},
		{IP: "127.0.0.1", Port: "9003"},
	}

	msg := node.NewMessage(
		route,
		"Hello with UUID",
	)

	data, _ := json.Marshal(msg)
	conn, _ := net.Dial("tcp", "127.0.0.1:9002") // Start from B
	defer conn.Close()
	conn.Write(data)

	time.Sleep(1 * time.Second)
}
