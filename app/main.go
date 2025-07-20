package main

import (
	"fmt"
	"starmesh/app/node"
	"starmesh/common/model"
	"starmesh/common/util"
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

	lat, lon, err := util.GetGeoLocation()
	if err != nil {
		fmt.Println("⚠️  위도/경도 조회 실패:", err)
		lat, lon = 0.0, 0.0
	}

	node.Bootstrap("http://localhost:8080", lat, lon)
	go node.StartServer()
	node.StartGossipLoop()
	node.StartPeerCleanupLoop()
	node.StartPeerPingLoop()

	// 서버 계속 실행 상태로 유지
	select {}
}
