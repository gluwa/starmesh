package main

import (
	"starmesh/app/node"
	"starmesh/common/model"
)

func main() {
	/*
			fmt.Println("▶ Running peer registration demo...")
			node.TestRegisterPeerDemo()

			fmt.Println("\n▶ Running connection demo...")
			node.TestConnectionDemo()

		fmt.Println("\n▶ Running N-Hop demo...")
		node.TestNHop()
	*/

	node := &node.Node{
		Port:  "9001",
		Peers: make(map[string]model.Peer),
	}

	go node.StartServer()
	node.StartGossipLoop()
	node.StartPeerCleanupLoop()
	node.StartPeerPingLoop()

	// 서버 계속 실행 상태로 유지
	select {}
}
