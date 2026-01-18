# GDRM (Go DynamoDB oRM)

DynamoDB Single Table Design을 위한 Go 라이브러리

[![Go Reference](https://pkg.go.dev/badge/github.com/zkfmapf123/gdrm.svg)](https://pkg.go.dev/github.com/zkfmapf123/gdrm)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## 특징

- 🚀 **심플한 API** - 직관적인 메서드로 DynamoDB 작업 수행
- 📦 **Single Table Design** - PK/SK 기반 설계 지원
- 🔄 **배치 처리** - 25개씩 자동 분할되는 배치 삽입
- 🎯 **제네릭 지원** - 타입 안전한 Marshal/Unmarshal
- 📝 **내장 로깅** - 디버깅을 위한 트레이스 로그

## 설치

```bash
go get github.com/zkfmapf123/gdrm
```

## 빠른 시작

```go
package main

import (
    "context"
    "log"

    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
    gdrm "github.com/zkfmapf123/gdrm"
)

type User struct {
    PK   string `dynamodbav:"PK"`
    SK   string `dynamodbav:"SK"`
    Name string `dynamodbav:"Name"`
    Age  int    `dynamodbav:"Age"`
}

func main() {
    // AWS 설정
    cfg, _ := config.LoadDefaultConfig(context.Background())
    dynamoClient := dynamodb.NewFromConfig(cfg)

    // 클라이언트 생성
    client := gdrm.NewDDB(dynamoClient)

    // 테이블 설정
    client.AddTable("my_table", gdrm.DDBTableParams{
        IsCreate:        true,
        IsPK:            true,
        PkAttributeType: types.ScalarAttributeTypeS,
        IsSK:            true,
        SkAttributeType: types.ScalarAttributeTypeS,
    })

    // 테이블 생성
    client.Start(true)

    // 데이터 삽입
    client.Insert("my_table", User{
        PK:   "USER#123",
        SK:   "#PROFILE",
        Name: "tom",
        Age:  32,
    })

    // 데이터 조회
    item, _ := client.FindByKey("my_table", "USER#123", "#PROFILE")
    user := gdrm.MarshalMap[User](item)

    log.Printf("Name: %s, Age: %d", user.Name, user.Age)
}
```

## API

### Client Functions

| 함수 | 설명 |
|------|------|
| `NewDDB(client)` | DynamoDB 클라이언트 생성 |
| `AddTable(name, params)` | 테이블 설정 추가 |
| `Start(isCreate)` | 테이블 생성 시작 |

### Insert Functions

| 함수 | 설명 |
|------|------|
| `Insert(tableName, item)` | 단건 삽입 (PK 중복 체크) |
| `InsertBatch(tableName, items)` | 배치 삽입 (25개씩 자동 분할) |

### Select Functions

| 함수 | 설명 |
|------|------|
| `FindByKey(tableName, pk, sk)` | PK/SK로 단건 조회 |
| `FindByKeyUseExpression(tableName, pk, sk, limit, params)` | Expression 조건부 조회 |

### Marshal Functions

| 함수 | 설명 |
|------|------|
| `MarshalMap[T](item)` | 단건 결과 타입 변환 |
| `MarshalMaps[T](items)` | 복수 결과 타입 변환 |

## 사용 예제

### 단건 조회

```go
item, err := client.FindByKey("my_table", "USER#123", "#PROFILE")
if err != nil {
    log.Fatal(err)
}

user := gdrm.MarshalMap[User](item)
```

### Expression을 사용한 조회

```go
items, err := client.FindByKeyUseExpression(
    "my_table",
    "TEAM#DEV",
    "",
    100,
    gdrm.RangeParams{
        KeyConditionExpression: "PK = :pk",
        ExpressionAttributeValues: map[string]types.AttributeValue{
            ":pk": &types.AttributeValueMemberS{Value: "TEAM#DEV"},
        },
    },
)

users := gdrm.MarshalMaps[User](items)
```

### begins_with 사용

```go
items, err := client.FindByKeyUseExpression(
    "my_table",
    "USER#123",
    "ORDER#",
    50,
    gdrm.RangeParams{
        KeyConditionExpression: "PK = :pk AND begins_with(SK, :sk)",
        ExpressionAttributeValues: map[string]types.AttributeValue{
            ":pk": &types.AttributeValueMemberS{Value: "USER#123"},
            ":sk": &types.AttributeValueMemberS{Value: "ORDER#"},
        },
    },
)
```

### 배치 삽입

```go
users := []any{
    User{PK: "USER#1", SK: "#PROFILE", Name: "tom", Age: 32},
    User{PK: "USER#2", SK: "#PROFILE", Name: "jane", Age: 28},
    User{PK: "USER#3", SK: "#PROFILE", Name: "mike", Age: 30},
}

err := client.InsertBatch("my_table", users)
```

## Single Table Design

DynamoDB Single Table Design의 핵심 원칙:

```
PK = 큰 그룹 (폴더)
SK = 세부 항목 (파일)
```

### 설계 예시

```go
// 유저 정보
{PK: "USER#1", SK: "#PROFILE", Name: "tom", Age: 32}

// 유저의 주문들
{PK: "USER#1", SK: "ORDER#001", ...}
{PK: "USER#1", SK: "ORDER#002", ...}

// 팀별 멤버 (양방향 조회를 위한 중복 저장)
{PK: "TEAM#DEV", SK: "USER#1", Name: "tom"}
```

자세한 설계 가이드는 [SINGLE_TABLE_ARCHITECTURE.md](./SINGLE_TABLE_ARCHITECTURE.md)를 참고하세요.

## Todo

- [ ] Backoff Limiter 추가 (Rate Limit)
- [ ] Update 함수 추가
- [ ] Delete 함수 추가
- [ ] GSI 지원
- [ ] Transaction 지원

## License

MIT License
