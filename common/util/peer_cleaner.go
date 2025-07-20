package util

import (
	"fmt"
	"sync"
	"time"

	"starmesh/common/model"
)

// StartPeerCleanupLoop starts a goroutine that removes peers not seen within expireDuration
func StartPeerCleanupLoop(peerMap *map[string]model.Peer, mux *sync.RWMutex, expireDuration time.Duration, label string) {
	go func() {
		for {
			time.Sleep(30 * time.Second)
			deadline := time.Now().Add(-expireDuration)

			mux.Lock()
			for k, peer := range *peerMap {
				if peer.LastSeen.Before(deadline) {
					delete(*peerMap, k)
					fmt.Printf("[%s Cleanup] Removed expired peer: %s\n", label, k)
				}
			}
			mux.Unlock()
		}
	}()
}
