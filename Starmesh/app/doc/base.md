# 🌐 Starmesh – P2P 기반 메시지 라우팅 시스템

**Starmesh**는 TOR-light 구조를 참고한 분산형 P2P 메시지 전송 시스템입니다.  
Seed Server 없이도 노드 간 직접 연결하여 메시지를 **여러 경유 노드를 통해 목적지까지 전달**합니다.  
실시간 시각화를 위한 Trace 기록, 장애 복구(재전송), 추후 암호화 확장을 고려한 구조로 설계되어 있습니다.

---

## 📁 프로젝트 구조

starmesh/  
├── go.mod  
├── main.go  
├── node/  
│ ├── node.go # Node 구조체, TCP 서버, Peer 등록 처리  
│ ├── peer.go # Peer 정의 및 연결 함수  
│ ├── message.go # 메시지 정의, 경유 전달, 재시도 처리  
│ ├── netutil.go # IP 정규화, 외부 IP 조회 유틸  
│ └── trace.go # 메시지 경로 Trace 기록(JSON 파일 저장)  
└── test/  
├── peer_test.go  
└── connection_test.go  


---

## 🧩 구현된 주요 기능

### ✅ Node 구조

- `Node`는 서버 역할(`Listen`)과 클라이언트 역할(`Dial`)을 동시에 수행
- 노드가 다른 노드로부터 메시지를 수신하면, **자동으로 Peer 등록**

```go
type Node struct {
    Port     string
    Peers    map[string]Peer
    PeersMux sync.RWMutex
}
```

✅ Peer 등록
 - IP:Port를 키로 하여 중복 없이 관리
 - 연결 요청이 들어오면 자동 등록
 - IPv6 루프백(::1) 주소 → 127.0.0.1로 정규화
 - GetOutboundIP()으로 실제 로컬 IP 확인

✅ Message 구조 및 라우팅
```go
type Message struct {
    MessageID string   // UUID
    Route     []Peer   // 다음 노드 리스트
    Payload   string   // 전달할 내용
    HopIndex  int      // 현재 위치
    Trace     []string // 경유 노드 목록
}
```
 - JSON 직렬화 후 TCP로 전송
 - Hop-by-hop 구조로 다음 노드에 전달
 - 마지막 노드에서는 Trace를 trace_log.json에 저장

✅ UUID 메시지 식별자
 - 메시지마다 UUID를 부여하여 Trace 분리
 - Viewer에서 메시지별 경로 구분 가능

✅ 메시지 전달 실패 시 재전송
 - 최대 3회 (MaxRetry)
 - 2초 간격 (RetryBackoff)
 - 실패 시 "FAILED:<ip:port>"를 Trace에 기록
 - Trace 파일에 실패 경로 포함하여 저장

✅ Trace 로그 파일 (trace_log.json)
한 메시지의 경로 기록 예:

```json
{
  "timestamp": "2025-07-17T22:30:10+09:00",
  "message_id": "fc62b3ab-4a1f-4f4f-b723-508038e92a60",
  "path": [
    "192.168.0.10:9001",
    "192.168.0.11:9002",
    "FAILED:192.168.0.12:9003"
  ],
  "payload": "Hello through the mesh"
}
```

## 🛰️ 향후 계획 (Todo)
 Seed Server 구현 (초기 Peer 정보 제공 및 등록)

 Gossip 기반 Peer 목록 공유

 실시간 Viewer 연동 (WebSocket 서버)

 Onion Routing 기반 암호화 구조

 Ping/Pong 기반 Peer 상태 체크

## 📦 의존성
 - Go ≥ 1.20
 - google/uuid

```bash
go get github.com/google/uuid
```

### ▶ 실행 예시
```bash
go run main.go
```
또는 테스트 실행:

```bash
go test -v ./test
```

## 👥 작성 및 기여자
 - 기획/설계: @wooghi
 - 개발: @wooghi