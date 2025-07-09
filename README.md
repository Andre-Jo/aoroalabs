# 이슈 관리 API

Go 언어로 구현된 이슈 관리 REST API 서버입니다.

## 실행 방법

### 1. 사전 준비
- Go 1.21 이상 설치 필요
- Git 설치 필요

### 2. 프로젝트 클론 및 실행
```bash
# 프로젝트 클론
git clone https://github.com/Andre-Jo/aoroalabs.git
cd aoroalabs

# 의존성 설치
go mod tidy

# 서버 실행
go run main.go
```

서버가 성공적으로 시작되면 포트 8080에서 실행됩니다.

## API 명세

### 1. 이슈 생성
- **URL**: `POST /issue`
- **설명**: 새로운 이슈를 생성합니다.
- **요청 예시**:
```json
{
    "title": "버그 수정 필요",
    "description": "로그인 페이지에서 오류 발생",
    "userId": 1
}
```
- **응답 예시** (201 Created):
```json
{
    "id": 1,
    "title": "버그 수정 필요",
    "description": "로그인 페이지에서 오류 발생",
    "status": "IN_PROGRESS",
    "user": {
        "id": 1,
        "name": "김개발"
    },
    "createdAt": "2025-06-02T10:00:00Z",
    "updatedAt": "2025-06-02T10:00:00Z"
}
```

### 2. 이슈 목록 조회
- **URL**: `GET /issues`
- **설명**: 이슈 목록을 조회합니다. 상태별 필터링 가능합니다.
- **쿼리 파라미터**:
  - `status` (선택): PENDING, IN_PROGRESS, COMPLETED, CANCELLED
- **응답 예시** (200 OK):
```json
{
    "issues": [
        {
            "id": 1,
            "title": "버그 수정 필요",
            "description": "로그인 페이지에서 오류 발생",
            "status": "PENDING",
            "createdAt": "2025-06-02T10:00:00Z",
            "updatedAt": "2025-06-02T10:05:00Z"
        }
    ]
}
```

### 3. 이슈 상세 조회
- **URL**: `GET /issue/:id`
- **설명**: 특정 이슈의 상세 정보를 조회합니다.
- **응답 예시** (200 OK):
```json
{
    "id": 1,
    "title": "버그 수정 필요",
    "description": "로그인 페이지에서 오류 발생",
    "status": "PENDING",
    "createdAt": "2025-06-02T10:00:00Z",
    "updatedAt": "2025-06-02T10:05:00Z"
}
```

### 4. 이슈 수정
- **URL**: `PATCH /issue/:id`
- **설명**: 기존 이슈를 수정합니다.
- **요청 예시**:
```json
{
    "title": "로그인 버그 수정",
    "status": "IN_PROGRESS",
    "userId": 2
}
```
- **응답 예시** (200 OK):
```json
{
    "id": 1,
    "title": "로그인 버그 수정",
    "description": "로그인 페이지에서 오류 발생",
    "status": "IN_PROGRESS",
    "user": {
        "id": 2,
        "name": "이디자인"
    },
    "createdAt": "2025-06-02T10:00:00Z",
    "updatedAt": "2025-06-02T10:10:00Z"
}
```

## API 테스트 방법

### curl을 사용한 테스트

#### 1. 이슈 생성 (담당자 있음)
```bash
curl -X POST http://localhost:8080/issue \
  -H "Content-Type: application/json" \
  -d '{
    "title": "버그 수정 필요",
    "description": "로그인 페이지에서 오류 발생",
    "userId": 1
  }'
```

#### 2. 이슈 생성 (담당자 없음)
```bash
curl -X POST http://localhost:8080/issue \
  -H "Content-Type: application/json" \
  -d '{
    "title": "새로운 기능 요청",
    "description": "사용자 대시보드 개선"
  }'
```

#### 3. 전체 이슈 목록 조회
```bash
curl -X GET http://localhost:8080/issues
```

#### 4. 상태별 이슈 조회
```bash
curl -X GET "http://localhost:8080/issues?status=PENDING"
```

#### 5. 이슈 상세 조회
```bash
curl -X GET http://localhost:8080/issue/1
```

#### 6. 이슈 수정
```bash
curl -X PATCH http://localhost:8080/issue/1 \
  -H "Content-Type: application/json" \
  -d '{
    "title": "수정된 제목",
    "status": "IN_PROGRESS",
    "userId": 2
  }'
```

#### 7. 담당자 제거 (userId를 0으로 설정)
```bash
curl -X PATCH http://localhost:8080/issue/1 \
  -H "Content-Type: application/json" \
  -d '{
    "userId": 0
  }'
```

### Postman을 사용한 테스트
1. Postman 실행
2. 새 Collection 생성
3. 위의 curl 명령어들을 Postman 요청으로 변환하여 테스트

## 기본 사용자 정보
시스템에 기본적으로 등록된 사용자들:
- ID: 1, 이름: "김개발"
- ID: 2, 이름: "이디자인"  
- ID: 3, 이름: "박기획"

## 비즈니스 규칙
1. **이슈 상태**: PENDING, IN_PROGRESS, COMPLETED, CANCELLED
2. **담당자 할당 규칙**:
   - 담당자가 있는 경우 → IN_PROGRESS
   - 담당자가 없는 경우 → PENDING
3. **수정 규칙**:
   - COMPLETED, CANCELLED 상태에서는 수정 불가
   - 담당자 없이 PENDING, CANCELLED 이외의 상태로 변경 불가
   - 담당자 제거 시 자동으로 PENDING 상태로 변경

## 에러 응답 형식
```json
{
    "error": "에러 메시지",
    "code": 400
}
```

## 프로젝트 구조
```
aoroalabs/
├── main.go              # 메인 애플리케이션
├── go.mod              # Go 모듈 설정
├── models/             # 데이터 모델
│   └── models.go
├── storage/            # 데이터 저장소
│   └── storage.go
├── handlers/           # HTTP 핸들러
│   └── handlers.go
└── README.md           # 프로젝트 문서
```