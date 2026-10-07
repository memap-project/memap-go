# Memap Go Client SDK

Офіційний Go SDK клієнт для високопродуктивного in-memory сховища [Memap](https://github.com/memap-project/memap).

Клієнт взаємодіє з сервером через протокол TCP з бінарною серіалізацією Protobuf (`protorw`), є повністю потокобезпечним (thread-safe) та підтримує модульні підклієнти для всіх структур даних.

## Встановлення

```bash
go get github.com/memap-project/memap-go
```

## Швидкий старт

```go
package main

import (
 "context"
 "fmt"
 "log"
 "time"

 "github.com/memap-project/memap-go"
)

func main() {
 ctx := context.Background()

 // Створення клієнта
 client, err := memap.New("localhost:8080",
  memap.WithTimeout(5*time.Second),
  memap.WithDialTimeout(3*time.Second),
 )
 if err != nil {
  log.Fatalf("failed to connect: %v", err)
 }
 defer client.Close()

 // Перевірка зв'язку
 pong, err := client.Ping(ctx)
 if err != nil {
  log.Fatalf("ping failed: %v", err)
 }
 fmt.Println("Ping:", pong) // PONG

	// 1. Робота через підклієнти (KV, Hash, Counter, RBuffer, Set)
 _ = client.Create(ctx, "users")
 kv := client.KV("users")
 _ = kv.Set(ctx, "user:1", "John Doe", 3600)
 val, _ := kv.Get(ctx, "user:1")
 fmt.Println("User:", val)

 // 2. Робота через простір імен (Namespace-scoped API)
 users := client.Namespace("users")
 _ = users.KV().Set(ctx, "user:2", "Jane Doe", 0)
 val2, _ := users.KV().Get(ctx, "user:2")
 fmt.Println("User 2:", val2)

 count, _ := users.Counter().IncrBy(ctx, "logins", 1)
 fmt.Println("Logins:", count)
}
```

## Модульні підклієнти

### Key-Value (`client.KV(namespace)`)

```go
kv := client.KV("app")

// Set & Get
err := kv.Set(ctx, "session:123", "data", 60) // TTL 60s
val, err := kv.Get(ctx, "session:123")

// TTL & Expire & Delete
ttl, err := kv.TTL(ctx, "session:123")
err = kv.Expire(ctx, "session:123", 120)
err = kv.Del(ctx, "session:123")
```

### Hash Map (`client.Hash(namespace)`)

```go
hash := client.Hash("app")

// Створення та робота з полями
_ = hash.Set(ctx, "profile:1", 3600)
_ = hash.FSet(ctx, "profile:1", "name", "Alex")
_ = hash.FSet(ctx, "profile:1", "city", "Kyiv")

name, _ := hash.FGet(ctx, "profile:1", "name")
allFields, _ := hash.Get(ctx, "profile:1") // map[string]string
keys, _ := hash.Keys(ctx, "profile:1")     // []string
vals, _ := hash.Values(ctx, "profile:1")   // []string
exists, _ := hash.Exists(ctx, "profile:1")
length, _ := hash.Len(ctx, "profile:1")
_ = hash.FDel(ctx, "profile:1", "city")
_ = hash.Del(ctx, "profile:1")
```

### Counter (`client.Counter(namespace)`)

```go
counter := client.Counter("metrics")

// Встановлення ліміту (ініціалізує лічильник, якщо він ще не існує)
_ = counter.SetLimit(ctx, "hits", 1000)
_ = counter.Expire(ctx, "hits", 3600)

newVal, _ := counter.IncrBy(ctx, "hits", 1)
newVal, _ = counter.DecrBy(ctx, "hits", 1)
count, _ := counter.Get(ctx, "hits")
limit, _ := counter.GetLimit(ctx, "hits")
_ = counter.SetLimit(ctx, "hits", 2000)
```

### RingBuffer (`client.RBuffer(namespace)`)

```go
buf := client.RBuffer("logs")

// Буфер місткістю 10 елементів
_ = buf.Init(ctx, "recent_events", 10, 0)
_ = buf.Push(ctx, "recent_events", "event_1")
_ = buf.Push(ctx, "recent_events", "event_2")

item, _ := buf.Pop(ctx, "recent_events") // вилучає найстаріший
items, _ := buf.Slice(ctx, "recent_events") // всі елементи
first, _ := buf.Peek(ctx, "recent_events")
last, _ := buf.Back(ctx, "recent_events")
length, _ := buf.Len(ctx, "recent_events")
capacity, _ := buf.Cap(ctx, "recent_events")
_ = buf.Reset(ctx, "recent_events")
```

### Set (`client.Set(namespace)`)

```go
set := client.Set("tags")

// Додавання та перевірка елементів
_ = set.Add(ctx, "article:1", "golang")
_ = set.Add(ctx, "article:1", "nosql")
isMember, _ := set.IsMember(ctx, "article:1", "golang") // true
count, _ := set.Card(ctx, "article:1")                   // 2
members, _ := set.Members(ctx, "article:1")              // []string{"golang", "nosql"}
_ = set.Remove(ctx, "article:1", "nosql")
_ = set.Expire(ctx, "article:1", 3600)
ttl, _ := set.TTL(ctx, "article:1")
```

### Namespace Management (`client`)

```go
_ = client.Create(ctx, "my_namespace")
_ = client.Drop(ctx, "my_namespace")
_ = client.Flush(ctx) // очищує дані в усіх namespace
_ = client.Erase(ctx) // видаляє всі namespace
```
