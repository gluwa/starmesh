package node

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"sync"
	"time"

	"starmesh/common/model"
	"starmesh/common/util"
)

type Node struct {
	Port     string
	Peers    map[string]model.Peer
	PeersMux sync.RWMutex
}

// StartServer starts a TCP server and accepts connections from other peers
func (n *Node) StartServer() {
	listener, err := net.Listen("tcp", ":"+n.Port)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	fmt.Println("[Server] Listening on port", n.Port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go n.handleConnection(conn)
	}
}

// handleConnection processes an incoming TCP connection
func (n *Node) handleConnection(conn net.Conn) {
	defer conn.Close()

	remoteAddr := conn.RemoteAddr().String()
	host, port, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		fmt.Println("[Server] Failed to parse remote address:", err)
		return
	}

	// Peer 등록
	n.RegisterPeer(model.Peer{
		IP:   util.NormalizeIP(host),
		Port: port,
		Lat:  0.0,
		Lon:  0.0,
	})

	// 메시지 읽기
	buffer := make([]byte, 4096)
	nBytes, _ := conn.Read(buffer)

	// Try to unmarshal as GossipMessage
	var gmsg GossipMessage
	if err := json.Unmarshal(buffer[:nBytes], &gmsg); err == nil && gmsg.SenderID != "" {
		n.HandleGossip(gmsg, conn)
		return
	}

	var msg Message
	err = json.Unmarshal(buffer[:nBytes], &msg)
	if err != nil {
		fmt.Printf("[Server] Received raw message: %s\n", string(buffer[:nBytes]))
		return
	}

	msg.HopIndex++

	selfID := fmt.Sprintf("%s:%s", util.GetOutboundIP(), n.Port)
	msg.Trace = append(msg.Trace, selfID)

	// 마지막 노드일 경우 로그 기록
	if msg.HopIndex >= len(msg.Route) {
		fmt.Println("[End] Message payload:", msg.Payload)
		SaveTraceToFile(&msg)
		return
	}

	msg.SendToNextHop()
}

// RegisterPeer adds or updates a peer
func (n *Node) RegisterPeer(p model.Peer) {
	p.IP = util.NormalizeIP(p.IP)
	key := p.IP + ":" + p.Port

	// 자기 자신이면 등록 안 함
	if p.IP == util.GetOutboundIP() && p.Port == n.Port {
		return
	}

	n.PeersMux.Lock()
	defer n.PeersMux.Unlock()

	// 중복이면 무시 (하지만 LastSeen은 갱신할 수도 있음)
	if _, exists := n.Peers[key]; exists {
		return
	}

	p.LastSeen = time.Now()
	n.Peers[key] = p

	fmt.Printf("[Server] RegisterPeer: %s at (%f, %f)\n", key, p.Lat, p.Lon)
}

func (n *Node) PingPeer(peer model.Peer) bool {
	address := peer.IP + ":" + peer.Port
	conn, err := net.DialTimeout("tcp", address, 1*time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func (n *Node) HandleGossip(gmsg GossipMessage, conn net.Conn) {
	count := 0
	for _, peer := range gmsg.Peers {
		n.RegisterPeer(peer)
		count++
	}

	fmt.Printf("[Gossip] Received %d peers from %s\n", count, gmsg.SenderID)

	// Pull 응답 전송
	if gmsg.Type == "push" && conn != nil {
		n.PeersMux.RLock()
		defer n.PeersMux.RUnlock()

		resp := GossipMessage{
			Type:     "pull",
			SenderID: fmt.Sprintf("%s:%s", util.GetOutboundIP(), n.Port),
			Peers:    util.SelectRandomPeers(n.Peers, 5),
		}
		respData, _ := json.Marshal(resp)
		conn.Write(respData)
	}
}

func (n *Node) StartGossipLoop() {
	go func() {
		for {
			time.Sleep(10 * time.Second)

			n.PeersMux.RLock()
			var peers []model.Peer
			for _, p := range n.Peers {
				peers = append(peers, p)
			}
			n.PeersMux.RUnlock()

			if len(peers) == 0 {
				continue
			}

			target := peers[rand.Intn(len(peers))]
			n.SendGossipTo(target)
		}
	}()
}

func (n *Node) StartPeerCleanupLoop() {
	util.StartPeerCleanupLoop(&n.Peers, &n.PeersMux, 10*time.Minute, "App")
}

func (n *Node) StartPeerPingLoop() {
	go func() {
		for {
			time.Sleep(20 * time.Second)

			n.PeersMux.Lock()
			for k, p := range n.Peers {
				if !n.PingPeer(p) {
					fmt.Printf("[Ping] Peer unreachable: %s\n", k)
					delete(n.Peers, k)
				}
			}
			n.PeersMux.Unlock()
		}
	}()
}
