package node

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"starmesh/common/model"
	"starmesh/common/util"
)

// Bootstrap connects to the Seed Server to register this node and fetch initial peer list
func (n *Node) Bootstrap(seedURL string, lat, lon float64) {
	// 1. 내 정보 생성
	self := model.Peer{
		IP:   util.GetOutboundIP(),
		Port: n.Port,
		Lat:  lat,
		Lon:  lon,
	}

	// 2. Seed Server에 등록
	body, _ := json.Marshal(self)
	resp, err := http.Post(seedURL+"/register", "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("[Bootstrap] Failed to register with Seed Server:", err)
		return
	}
	resp.Body.Close()
	fmt.Printf("[Bootstrap] Registered with Seed Server: %s:%s\n", self.IP, self.Port)

	// 3. Peer 목록 가져오기
	res, err := http.Get(seedURL + "/peers")
	if err != nil {
		fmt.Println("[Bootstrap] Failed to fetch peer list:", err)
		return
	}
	defer res.Body.Close()

	respBody, _ := io.ReadAll(res.Body)

	var peers []model.Peer
	if err := json.Unmarshal(respBody, &peers); err != nil {
		fmt.Println("[Bootstrap] Failed to parse peer list:", err)
		return
	}

	count := 0
	for _, p := range peers {
		n.RegisterPeer(p)
		count++
	}
	fmt.Printf("[Bootstrap] Loaded %d peers from Seed Server\n", count)
}
