**Table Of Contents**
- [RORM (Raw Query ORM) Description](#RORM-Raw-Query-ORM-Description)
- [Benchmarking vs XORM](#Benchmarking-vs-XORM)
- [Support Database](#Support-Database)
- [Installation](#Installation)
- [Features](#Features)
- [How To Use](#How-To-Use)
  - [Configure the Host](#Configure-the-Host)
  - [Init New Engine](#Init-New-Engine)
  - [Init the models](#Init-the-models)
  - [Create New SQL Select Query](#Create-New-SQL-Select-Query)
    - [Get All Data](#Get-All-Data)
      - [With SQL Raw](#With-SQL-Raw)
      - [WITH Query Builder](#WITH-Query-Builder)
    - [Get Multiple Result Data with Where Condition](#Get-Multiple-Result-Data-with-Where-Condition)
    - [Get Single Result Data with Where Condition](#Get-Single-Result-Data-with-Where-Condition)
  - [Create, Update, Delete Query](#Create-Update-Delete-Query)
    - [Insert](#Insert)
      - [Single Insert](#Single-Insert)
      - [Multiple Insert](#Multiple-Insert)
    - [Update](#Update)
    - [Delete](#Delete)
# RORM (Raw Query ORM) Description
Raw Query ORM Library for golang (postgres, mysql)
other RDBMS coming soon

NoSQL query will be coming soon too


# Benchmarking vs XORM
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

# Support Database
| No   | Database   |
| :--- | :--------- |
| 1    | MySQL      |
| 2    | Postgres   |
| 3    | SQL Server |

# Installation
```go
go get github.com/radityaapratamaa/rorm
```

import to your project

```go
import "github.com/radityaapratamaa/rorm"
```
# Features
| Read          | CUD    |
| :------------ | :----- |
| Select        | Insert |
| SelectSum     | Update |
| SelectAverage | Delete |
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

# How To Use
## Configure the Host
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
## Init New Engine
```go
// Please make sure the variable "dbConfig" pass by reference (pointer)
db, err := rorm.New(dbConfig)
if err != nil {
    log.Fatalln("Cannot Connet to database")
}
log.Println("Success Connect to Database")
```

## Init the models
```sql
    -- We Have a table with this structure
    CREATE TABLE Student (
        id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
        name VARCHAR(100) NULL,
        address TEXT NULL,
        is_active BOOLEAN NULL,
        birth_date DATE NULL, 
    )
```

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

## Create New SQL Select Query
### Get All Data
#### With SQL Raw
```go
    if err := db.SQLRaw(`SELECT name, address, birth_date FROM Student [JOIN ..... ON .....] 
    [WHERE ......... [ORDER BY ...] [LIMIT ...]]`).Get(&studentList); err != nil {
        log.Fatalln(err.Error())
    }
    log.Println("result is, ", studentList)
```
#### WITH Query Builder
```go
    // Get All Students data
    if err := db.Get(&studentList); err != nil {
        log.Fatalln(err.Error())
    }
    log.Println("result is, ", studentList)
```
```sql
    --  it will generate : 
    SELECT * FROM student
```
### Get Multiple Result Data with Where Condition
```go
    // Get Specific Data
    if err := db.Select("name, address, birth_date").Where("is_active", 1).Get(&studentList); err != nil {
        log.Fatalln(err.Error())
    }
    log.Println("result is, ", studentList)
```
```sql
-- it will generate: (prepared Statement)
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
-- it will generate: 
SELECT name, address, birth_date FROM student WHERE is_active = ? AND name LIKE ?
```
### Get Single Result Data with Where Condition
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

## Create, Update, Delete Query
### Insert
#### Single Insert
```go
    dtStudent := Student{
        Name: "test",
        Address: "test",
        IsActive: 1,
        BirthDate: "2010-01-01",
    }

    affected, err := db.Insert(&dtStudent)
    if err != nil {
        log.Fatalln("Error When Insert")
    }

    if affected > 0 {
        log.Println("Success Insert")
    }
```
```sql
    -- it will generate : (mysql)
    INSERT INTO Student (name, address, is_active, birth_date) VALUES (?,?,?,?)
    -- prepared Values :
    -- ('test', 'test', 1, '2010-01-01')
```
#### Multiple Insert
```go
    dtStudents := []Student{
        Student{
            Name: "test",
            Address: "test",
            IsActive: 1,
            BirthDate: "2010-01-01",
        },
        Student{
            Name: "test2",
            Address: "test2",
            IsActive: 1,
            BirthDate: "2010-01-02",
        },
        Student{
            Name: "test3",
            Address: "test3",
            IsActive: 1,
            BirthDate: "2010-01-03",
        },
    }

    affected, err := db.Insert(&dtStudent)
    if err != nil {
        log.Fatalln("Error When Insert")
    }

    if affected > 0 {
        log.Println("Success Insert")
    }
```
```sql
    -- it will generate : (mysql)
    INSERT INTO Student (name, address, is_active, birth_date) VALUES (?,?,?,?)
    -- prepared Values :
    -- 1. ('test', 'test', 1, '2010-01-01')
    -- 2. ('test2', 'test2', 1, '2010-01-02')
    -- 3. ('test3', 'test3', 1, '2010-01-03')
```

### Update
```go
    dtStudent := Student{
        Name: "change",
        Address: "change",
        IsActive: 1,
        BirthDate: "2010-01-10",
    }

    affected, err := db.Where("id", 1).Update(&dtStudent)
    if err != nil {
        log.Fatalln("Error When Update")
    }

    if affected > 0 {
        log.Println("Success Update")
    }
```
```sql
    -- it will generate : (mysql)
    UPDATE Student SET name = ?, address = ?, is_active = ?, birth_date = ? WHERE id = ?
    -- prepared Values :
    -- ('change', 'change', 1, '2010-01-10', 1)
```
### Delete
```go
    affected, err := db.Where("id", 1).Delete(&Student{})
    if err != nil {
        log.Fatalln("Error When Delete")
    }

    if affected > 0 {
        log.Println("Success Delete")
    }
```
```sql
    -- it will generate : (mysql)
    DELETE FROM Student WHERE id = ?
    -- prepared Values :
    -- (1)
```


