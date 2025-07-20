
# 🌱 Starmesh Seed Server

**Seed Server**는 Starmesh P2P 네트워크에서 새로운 노드가 진입할 때 최초 Peer 정보를 받아갈 수 있도록 도와주는 경량화 HTTP 서버입니다.  
네트워크 확산은 Gossip 알고리즘에 의해 이루어지며, Seed Server는 초기에만 사용됩니다.

---

## ✅ 주요 기능

| 기능 | 설명 |
|------|------|
| `POST /register` | 새로운 Peer를 등록 |
| `GET /peers` | 등록된 Peer 중 일부를 반환 |
| 자동 Peer 정리 | 마지막 확인 시각(`lastSeen`) 기준으로 오래된 Peer 자동 제거 |
| 공통 모듈 공유 | `common/model`, `common/util` 을 통해 app과 구조/로직 공유 |

---

## 📦 구조

```bash
seedserver/
├── main.go                  # HTTP 서버 시작점
├── handler/
│   └── peer_handler.go      # API 핸들러 (등록 / 조회)
├── storage/
│   └── peer_store.go        # Peer 저장 및 동기화
```

### 공통 모듈:

```bash
common/
├── model/
│   └── peer.go              # Peer 구조체 정의
└── util/
    ├── netutil.go           # NormalizeIP, GetOutboundIP 등
    └── peer_cleaner.go      # 오래된 Peer 제거 루프 (app, seedserver 공용)
```

## 🧩 Peer 구조체 (common/model/peer.go)
```go
type Peer struct {
    IP       string    `json:"ip"`
    Port     string    `json:"port"`
    Lat      float64   `json:"lat"`
    Lon      float64   `json:"lon"`
    LastSeen time.Time `json:"lastSeen"`
}
```

## 📡 API 명세
### POST /register
신규 Peer 등록

```http
POST /register
Content-Type: application/json

{
  "ip": "192.168.0.10",
  "port": "9001",
  "lat": 37.5665,
  "lon": 126.9780
}
```
응답: 200 OK

GET /peers
일부 Peer 목록을 반환 (최대 10개)

```bash
GET /peers
```
응답 예시:
```json
[
  {
    "ip": "192.168.0.10",
    "port": "9001",
    "lat": 37.5665,
    "lon": 126.978,
    "lastSeen": "2025-07-20T21:56:05.736516+09:00"
  }
]
```

## 🔄 오래된 Peer 정리
 - common/util/peer_cleaner.go에 정의된 함수 사용

```go
util.StartPeerCleanupLoop(&storage.PeerMap, &storage.Mux, 10*time.Minute, "Seed")
```
 - 30초마다 루프
 - LastSeen 기준 10분 이상 경과된 Peer는 제거
 - App과 Seed Server가 동일 함수 공유

## 🚀 실행 방법
```bash
go run main.go
```
테스트:

```bash
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"ip":"192.168.0.10", "port":"9001", "lat":37.5665, "lon":126.9780}'

curl http://localhost:8080/peers
```

## 🧭 향후 확장 가능성
- /status: 전체 Peer 수, 서버 시간 등 확인 API
- IP 중복/자기자신 등록 필터
- Peer Score, TTL, 거리 기반 필터 등 확장
- GCP Cloud Run 배포

## 👥 개발자 참고
- Seed Server는 Peer 확산의 Entry Point일 뿐,
Gossip 전파, 라우팅은 App 노드 간 직접 수행됨.
- 전체 시스템은 common/ 모듈 기반으로 기능을 분산 유지함.