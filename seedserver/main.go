package main

import (
	"log"
	"net/http"
	"time"

	"starmesh/common/util"
	"starmesh/seedserver/handler"
	"starmesh/seedserver/storage"
)

func main() {
	// 먼저 백그라운드로 cleanup 루프 시작
	util.StartPeerCleanupLoop(&storage.PeerMap, &storage.Mux, 10*time.Minute, "Seed")

	// HTTP 핸들러 등록
	http.HandleFunc("/register", handler.RegisterPeerHandler)
	http.HandleFunc("/peers", handler.GetPeersHandler)

	log.Println("Seed Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil)) // 여기서 서버는 블로킹됨
}
