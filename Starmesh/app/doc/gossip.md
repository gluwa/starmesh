## 🔄 Gossip 기반 Peer 확산 기능

Starmesh는 Tor-like 구조의 노드 간 P2P 통신을 기반으로 하며, Gossip 알고리즘을 통해 Peer 목록을 주기적으로 확산합니다.  
Seed Server가 없어도 새로운 노드들이 네트워크에 동적으로 참여할 수 있도록 설계되었습니다.

---

### 📦 GossipMessage 구조

```go
type GossipMessage struct {
    Type     string   // "push" or "pull"
    SenderID string   // 보낸 노드의 IP:Port
    Peers    []Peer   // 전파하려는 Peer 목록
}
```

- Push: 노드가 자신이 알고 있는 Peer 목록 일부를 다른 노드에게 전송
- Pull: Push를 받은 노드가 응답으로 자신의 Peer 목록 일부를 다시 전송 (양방향 확산)

### ✅ 구현된 기능 요약
| 기능                   | 설명                                     |
| -------------------- | -------------------------------------- |
| 🎲 **랜덤 Peer 전파**    | 주기적으로 Peer 목록 중 일부를 무작위로 다른 Peer에게 전송  |
| ♻️ **Push + Pull**   | Gossip 전송 후 응답으로 Peer 목록을 받아 더욱 빠르게 확산 |
| 🧍 **자기 자신 등록 방지**   | 내 IP\:Port와 동일한 Peer는 등록하지 않음          |
| 🔁 **중복 Peer 등록 방지** | 이미 등록된 Peer는 다시 등록하지 않음                |
| ⏳ **LastSeen 필드 관리** | 각 Peer가 마지막으로 확인된 시각 기록 (`time.Time`)  |
| 🧹 **오래된 Peer 정리**   | 일정 시간(예: 10분) 이상 응답 없는 Peer는 제거        |
| 📡 **Ping 기반 연결 확인** | 주기적으로 Peer에 Ping → 실패 시 Peer 제거        |



#### 📤 Gossip 전송 로직 (SendGossipTo)
- Peer 목록 중 무작위로 N개 선택
- JSON으로 직렬화하여 TCP로 전송
- 전송 시 Type: "push"로 지정

#### 📥 Gossip 수신 로직 (HandleGossip)
- 수신한 Peer들을 필터링 후 등록
- 중복 제거 및 LastSeen 갱신
- Type == "push"일 경우, pull 응답 전송

#### ⏱ 주기적 동작 루프
- StartGossipLoop()
→ 10초마다 랜덤 Peer 1명에게 Gossip 전송
- StartPeerCleanupLoop()
→ 30초마다 LastSeen이 오래된 Peer 제거
- StartPeerPingLoop()
→ 20초마다 모든 Peer에 Ping → 실패 시 제거

#### 🔧 랜덤 Peer 선택 유틸
```go
func selectRandomPeers(peerMap map[string]Peer, maxCount int) []Peer
```
- Go 1.20 이상: rand.Shuffle 사용, 별도 rand.Seed 필요 없음
- 중복 없이 최대 N개의 Peer를 무작위로 선택

### 📁 관련 주요 파일
| 파일           | 역할                             |
| ------------ | ------------------------------ |
| `gossip.go`  | Gossip 메시지 전송/수신, Push+Pull 로직 |
| `node.go`    | Peer 등록, Ping 처리               |
| `netutil.go` | IP 정규화, 랜덤 Peer 선택 등 유틸        |
| `peer.go`    | Peer 구조체 정의                    |
| `message.go` | 일반 메시지와 함께 Trace 전달            |


### 🧠 향후 개선 아이디어
- Peer 신뢰도(성공률 기반 score) 관리
- Peer TTL과 Soft Eviction 정책 도입
- Gossip message routing 로그 시각화
- Dead peer 복구 재시도 전략

  
  
이 구현을 통해 Starmesh는 Seed Server 없이도 Peer를 자동 발견하고 확산할 수 있으며,
네트워크가 살아 숨 쉬는 분산 구조로 점진적으로 진화할 수 있는 기반을 마련합니다.