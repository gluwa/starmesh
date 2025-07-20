package storage

import (
	"starmesh/common/model"
	"sync"
	"time"
)

var (
	PeerMap = make(map[string]model.Peer) // 대문자로 export
	Mux     = sync.RWMutex{}
)

func AddPeer(p model.Peer) {
	key := p.IP + ":" + p.Port
	Mux.Lock()
	defer Mux.Unlock()
	p.LastSeen = time.Now()
	PeerMap[key] = p
}

func GetPeers(limit int) []model.Peer {
	Mux.RLock()
	defer Mux.RUnlock()

	peers := make([]model.Peer, 0, len(PeerMap))
	for _, p := range PeerMap {
		peers = append(peers, p)
	}

	if len(peers) > limit {
		return peers[:limit]
	}
	return peers
}
