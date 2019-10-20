# RORM (Raw Query ORM)
Raw Query ORM Library for golang (postgres, mysql)
other RDBMS coming soon

NoSQL query will be coming soon too

## Installation
```go
go get github.com/radityaapratamaa/rorm
```

import to your project

```go
import "github.com/radityaapratamaa/rorm"
```
## Features
| Read          | CUD    |
| :------------ | :----- |
| Select        | Insert |
| SelectSum     | Update |
| SelectAverage | Delete       |
| SelectMax     |        |
| SelectMin     |        |
| SelectCount   |        |
| Where         |        |
| WhereIn       |        |
| WhereNotIn    |        |
| WhereLike     |        |
| Or            |        |
| OrIn          |        |
| OrNotIn       |        |
| OrLike        |        |
| GroupBy       |        |
| Join          |        |
| Limit         |        |
| OrderBy       |        |
| Asc           |        |
| Desc          |        |

## How To Use
### Configure the Host
```go
// Mandatory DBConfig
dbConfig := &rorm.DbConfig{
    Host: "localhost",
    Driver: "(mysql | postgres | sqlserver)",
    Username: "your_username",
    DbName: "your_dbName",
}
```

All Property DBConfig
```go
dbConfig := &rorm.DbConfig{
    Host: "your_host", //mandatory
    Driver: "DB Driver", //mandatory
    Username: "dbUsername", //mandatory
    DbName:"database_name", //mandatory
    Password: "dbPass",
    DbScheme: "db Scheme", // for postgres scheme, default is "Public" Scheme
    Port: "port", //default 3306 (mysql), 5432 (postgres)
    Protocol: "db Protocol", //default is "tcp"
    DbInstance: "db Instance", // for sqlserver driver if necessary
}
```
### Init New Engine
```go
// Please make sure the variable "dbConfig" pass by reference (pointer)
db, err := rorm.New(dbConfig)
if err != nil {
    log.Fatalln("Cannot Connet to database")
}
log.Println("Success Connect to Database")
```

### Create New SQL Select Query
#### Init the models
```go
    // Init the models (struct name MUST BE SAME with table name)
    type Student struct {
        // json tag, filled based on column name
        Name string `json:"name"`
        Address string `json:"address"`
        // IsActive is boolean field
        IsActive int `json:"is_active"`
        // BirthDate is "Date" data type column field
        BirthDate string `json:"birth_date"`
    }

    // init student struct to variable
    var studentList []Student
    var student Student
```
#### Get All Data
```go
    // Get All Students data
    if err := db.Get(&studentList); err != nil {
        log.Fatalln(err.Error())
    }
    log.Println("result is, ", studentList)
```
```sql
    // it will generate : 
    SELECT * FROM student
```
#### Get Multiple Result Data with Where Condition
```go
    // Get Specific Data
    if err := db.Select("name, address, birth_date").Where("is_active", 1).Get(&studentList); err != nil {
        log.Fatalln(err.Error())
    }
    log.Println("result is, ", studentList)
```
```sql
// it will generate: (prepared Statement)
SELECT name, address, birth_date FROM student WHERE is_active = ?
```
---
```go
    // Get Specific Data (other example)
    if err := db.Select("name", "address", "birth_date").
    Where("is_active", 1).WhereLike("name", "%Lorem%").
    Get(&studentList); err != nil {
        log.Fatalln(err.Error())
    }
    log.Println("result is, ", studentList)
```
```sql
// it will generate: 
SELECT name, address, birth_date FROM student WHERE is_active = ? AND name LIKE ?
```
#### Get Single Result Data with Where Condition
```go
    // Get Specific Data (single Result)
    if err := db.Select("name, address, birth_date").Where("id", 1).Get(&student); err != nil {
        log.Fatalln(err.Error())
    }
    log.Println("result is, ", student)
```
```sql
-- it will generate: 
SELECT name, address, birth_date FROM student WHERE id = ?
```

## Benchmarking vs XORM
**go test -bench=.**
```bash
goos: darwin
goarch: amd64
pkg: sample_file/test_rorm
BenchmarkInsertRorm-4               1000           1796572 ns/op             433 B/op          6 allocs/op
BenchmarkInsertRorm100-4            1000           1601391 ns/op             433 B/op          6 allocs/op
BenchmarkInsertRorm1000-4           1000           1531493 ns/op             436 B/op          6 allocs/op
BenchmarkInsertRorm10000-4          1000           1538203 ns/op             434 B/op          6 allocs/op
BenchmarkInsertXorm-4                500           2379465 ns/op            2320 B/op         59 allocs/op
BenchmarkInsertXorm100-4            1000           2528573 ns/op            2241 B/op         59 allocs/op
BenchmarkInsertXorm1000-4           1000           2245033 ns/op            2239 B/op         59 allocs/op
BenchmarkInsertXorm10000-4          1000           2275546 ns/op            2243 B/op         59 allocs/op
PASS
ok      sample_file/test_rorm   16.279s
```

