# Mongo Phoenix

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/logocomune/phoenix)
![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/logocomune/phoenix/go.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/logocomune/phoenix)](https://pkg.go.dev/github.com/logocomune/phoenix)
[![codecov](https://codecov.io/gh/logocomune/phoenix/graph/badge.svg?token=GGN3PHjyZV)](https://codecov.io/gh/logocomune/phoenix)
[![Go Report Card](https://goreportcard.com/badge/github.com/logocomune/phoenix)](https://goreportcard.com/report/github.com/logocomune/phoenix)

Mongo Phoenix is a lightweight, generic Go library that simplifies CRUD operations on MongoDB.
It wraps the official MongoDB Go driver with a clean, type-safe API and reduces boilerplate
by managing context timeouts and collection references automatically.

## Features

- **Generic CRUD** — type-safe `FindOne`, `FindAll`, `UpdateOne`, `UpdateMany`, `DeleteAll`, and `Count` methods.
- **Query builder** — compose filter documents from reusable `Option` functions (`ByKeyValue`, `InKeyValue`, `GTKeyValue`, etc.).
- **Update builder** — compose update documents from reusable `Option` functions (`Set`, `Inc`, `Push`, `Pull`, `AddToSet`, `SetOnInsert`).
- **Sort builder** — compose sort documents with `WithSort` / `WithSorts`.
- **Automatic timeout** — every operation applies the configured `time.Duration` as a context deadline.

## Installation

```bash
go get github.com/logocomune/phoenix
```

## Quick start

```go
package main

import (
    "context"
    "log"
    "time"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"

    phoenix "github.com/logocomune/phoenix"
)

type User struct {
    Name  string `bson:"name"`
    Email string `bson:"email"`
    Age   int    `bson:"age"`
}

func main() {
    client, _ := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
    db := client.Database("mydb")

    // Create a typed Phoenix wrapper bound to the "users" collection with a 5-second timeout.
    users := phoenix.New[User](db, "users", 5*time.Second)

    // FindOne
    user, found, err := users.FindOne(context.Background(), bson.M{"email": "alice@example.com"})
    if err != nil {
        log.Fatal(err)
    }
    if found {
        log.Printf("Found: %s", user.Name)
    }
}
```

---

## Package `phoenix` — CRUD operations

All methods share the same timeout configured in `New`.

### FindOne

Retrieves the first matching document. Returns `(doc, true, nil)` when found,
`(zero, false, nil)` when no document matches, or `(zero, false, err)` on error.

```go
users := phoenix.New[User](db, "users", 5*time.Second)

user, found, err := users.FindOne(ctx, bson.M{"name": "Alice"})
if err != nil {
    log.Fatal(err)
}
if found {
    fmt.Println(user.Email)
}
```

### FindAll

Retrieves all matching documents.

```go
activeUsers, err := users.FindAll(ctx, bson.M{"active": true})
if err != nil {
    log.Fatal(err)
}
for _, u := range activeUsers {
    fmt.Println(u.Name)
}
```

### Count

Counts documents matching the filter.

```go
n, err := users.Count(ctx, bson.M{"active": true})
fmt.Printf("%d active users\n", n)
```

### DeleteAll

Removes all matching documents and returns the deleted count.

```go
deleted, err := users.DeleteAll(ctx, bson.M{"active": false})
fmt.Printf("Removed %d inactive users\n", deleted)
```

### UpdateMany / UpdateOne

Returns `(matched, modified, error)`.

```go
matched, modified, err := users.UpdateMany(
    ctx,
    bson.M{"active": false},
    bson.M{"$set": bson.M{"archived": true}},
)
fmt.Printf("Matched: %d, Modified: %d\n", matched, modified)
```

---

## Package `query` — filter builder

Build complex MongoDB filter documents with composable helpers.

```go
import "github.com/logocomune/phoenix/query"

// Exact match + $in
filter := query.Generate(
    func(m bson.M) { query.ByKeyValue(m, "active", true) },
    func(m bson.M) { query.InKeyValue(m, "role", []string{"admin", "moderator"}) },
)
// → {"active": true, "role": {"$in": ["admin", "moderator"]}}

// Range filter: 18 ≤ age < 65
ageFilter := query.Generate(
    func(m bson.M) { query.GTEKeyValue(m, "age", 18) },
    func(m bson.M) { query.LTKeyValue(m, "age", 65) },
)
// → {"age": {"$gte": 18, "$lt": 65}}
```

Available helpers:

| Function        | MongoDB operator | Description                     |
|-----------------|------------------|---------------------------------|
| `ByKeyValue`    | —                | Exact equality match            |
| `InKeyValue`    | `$in`            | Field value is in the given set |
| `GTKeyValue`    | `$gt`            | Greater than                    |
| `GTEKeyValue`   | `$gte`           | Greater than or equal           |
| `LTKeyValue`    | `$lt`            | Less than                       |
| `LTEKeyValue`   | `$lte`           | Less than or equal              |

---

## Package `update` — update document builder

Build MongoDB update documents with composable operator helpers.

```go
import "github.com/logocomune/phoenix/update"

upd := update.Generate(
    update.Set("lastSeen", time.Now()),
    update.Inc("loginCount", 1),
    update.AddToSet("tags", "returning"),
)
// → {"$set": {"lastSeen": ...}, "$inc": {"loginCount": 1}, "$addToSet": {"tags": "returning"}}

matched, modified, err := users.UpdateMany(ctx, bson.M{"active": true}, upd)
```

Available operators:

| Function      | MongoDB operator | Description                                         |
|---------------|------------------|-----------------------------------------------------|
| `Set`         | `$set`           | Set a field value                                   |
| `Inc`         | `$inc`           | Increment a numeric field                           |
| `Push`        | `$push`          | Append a value to an array field                    |
| `Pull`        | `$pull`          | Remove matching values from an array field          |
| `AddToSet`    | `$addToSet`      | Add a value to a set (no duplicates)                |
| `SetOnInsert` | `$setOnInsert`   | Set a field only when an upsert inserts a document  |

---

## Package `sort` — sort document builder

Build MongoDB sort documents with composable helpers.

```go
import "github.com/logocomune/phoenix/sort"
import "go.mongodb.org/mongo-driver/bson"
import "go.mongodb.org/mongo-driver/mongo/options"

sortDoc := sort.Generate(
    sort.WithSort(bson.E{Key: "lastName", Value: 1}),   // ascending
    sort.WithSort(bson.E{Key: "createdAt", Value: -1}), // descending
)
// → [{"lastName": 1}, {"createdAt": -1}]

opts := options.Find().SetSort(sortDoc)
results, err := users.FindAll(ctx, bson.M{}, opts)
```

---

## Standalone functions

Every method on `Phoenix` has a corresponding package-level function that accepts
`*mongo.Database`, `collectionName`, and `timeout` directly — useful when you do
not want to create a `Phoenix` instance.

```go
user, found, err := phoenix.FindOne[User](ctx, 5*time.Second, db, "users", bson.M{"name": "Alice"})
count, err := phoenix.Count(ctx, 5*time.Second, db, "users", bson.M{"active": true})
```

## License

This project is licensed under the MIT License.
