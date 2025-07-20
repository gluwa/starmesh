package test

import (
	"fmt"
	"time"

	"starmesh/app/node"
	"testing"
)

func TestConnection(t *testing.T) {
	// Node A (서버 역할)
	go func() {
		serverNode := node.Node{Port: "8081"}
		serverNode.StartServer()
	}()

	// 서버가 먼저 실행될 수 있도록 잠시 대기
	time.Sleep(1 * time.Second)

	// Node B (클라이언트 역할)
	fmt.Println("[Client] Connecting to Node A...")
	node.ConnectToPeer("localhost:8081", "Hello from Node B")
}
