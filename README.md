**Table Of Contents**
- [RORM (Raw Query ORM) Description](#rorm-raw-query-orm-description)
- [Benchmarking vs XORM](#benchmarking-vs-xorm)
- [Support Database](#support-database)
- [Installation](#installation)
- [Features](#features)
- [How To Use](#how-to-use)
  - [Configure the Host](#configure-the-host)
  - [Init New Engine](#init-new-engine)
  - [Init the models](#init-the-models)
  - [Create New SQL Select Query](#create-new-sql-select-query)
    - [Get All Data](#get-all-data)
      - [With SQL Raw](#with-sql-raw)
      - [WITH Query Builder](#with-query-builder)
    - [Get Multiple Result Data with Where Condition](#get-multiple-result-data-with-where-condition)
    - [Get Single Result Data with Where Condition](#get-single-result-data-with-where-condition)
  - [Create, Update, Delete Query](#create-update-delete-query)
    - [Insert](#insert)
      - [Single Insert](#single-insert)
      - [Multiple Insert](#multiple-insert)
    - [Update](#update)
    - [Delete](#delete)
# RORM (Raw Query ORM) Description
Raw Query ORM is a Query Builder as light as raw query and as easy as ORM

# Benchmarking vs XORM
source : https://github.com/kihamo/orm-benchmark

command : ``` orm-benchmark -orm=xorm,rorm (-multi=1 default) ```

```bash
Reports: 

  2000 times - Insert
       raw:     3.32s      1660658 ns/op     568 B/op     14 allocs/op
      xorm:     5.38s      2688668 ns/op    2585 B/op     69 allocs/op
      rorm:     5.58s      2787619 ns/op    1052 B/op     14 allocs/op

   500 times - MultiInsert 100 row
       raw:     2.48s      4961667 ns/op  110997 B/op   1110 allocs/op
      xorm:     2.79s      5578304 ns/op  230925 B/op   4964 allocs/op
      rorm:    90.84s    181683062 ns/op   53478 B/op    709 allocs/op

  2000 times - Update
       raw:     1.66s       830352 ns/op     632 B/op     16 allocs/op
      xorm:     3.29s      1642816 ns/op    2914 B/op    108 allocs/op
      rorm:     3.40s      1700092 ns/op   13935 B/op    188 allocs/op

  4000 times - Read
      rorm:     3.35s       837601 ns/op    6213 B/op     85 allocs/op
       raw:     3.37s       841311 ns/op    1432 B/op     37 allocs/op
      xorm:     7.08s      1770239 ns/op    9762 B/op    268 allocs/op

  2000 times - MultiRead limit 100
       raw:     2.65s      1326544 ns/op   34720 B/op   1320 allocs/op
      rorm:     4.28s      2139107 ns/op   48668 B/op   1622 allocs/op
      xorm:     5.85s      2926974 ns/op  178591 B/op   7892 allocs/op
```

command: ``` orm-benchmark -orm=xorm,rorm,raw -multi=10 ```

```
Reports: 

 20000 times - Insert
       raw:    25.43s      1271440 ns/op     568 B/op     14 allocs/op
      rorm:    37.36s      1867861 ns/op    1027 B/op     14 allocs/op
      xorm:    39.14s      1956955 ns/op    2578 B/op     69 allocs/op

  5000 times - MultiInsert 100 row
       raw:    20.76s      4151136 ns/op  110910 B/op   1110 allocs/op
      xorm:    23.63s      4726541 ns/op  230813 B/op   4964 allocs/op
      rorm:    725.63s    145125198 ns/op   53322 B/op    707 allocs/op

 20000 times - Update
       raw:    11.96s       598057 ns/op     632 B/op     16 allocs/op
      xorm:    23.55s      1177686 ns/op    2914 B/op    108 allocs/op
      rorm:    33.16s      1658050 ns/op   13901 B/op    188 allocs/op

 40000 times - Read
       raw:    23.60s       589875 ns/op    1432 B/op     37 allocs/op
      rorm:    32.00s       799907 ns/op    6188 B/op     85 allocs/op
      xorm:    53.26s      1331431 ns/op    9762 B/op    268 allocs/op

 20000 times - MultiRead limit 100
       raw:    17.87s       893332 ns/op   34720 B/op   1320 allocs/op
      rorm:    41.47s      2073701 ns/op   48612 B/op   1622 allocs/op
      xorm:    49.59s      2479707 ns/op  178586 B/op   7892 allocs/op
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
| Feature       | Using                                               | Description                                                                                                      |
| :------------ | :-------------------------------------------------- | :--------------------------------------------------------------------------------------------------------------- |
| Select        | Select(cols ...string)                              | Specify the column will be query                                                                                 |
| SelectSum     | SelectSumn(col string)                              | Specify the single column to be summarize                                                                        |
| SelectAverage | SelectAverage(col string)                           | Specify the single column to average the value                                                                   |
| SelectMax     | SelectMax(col string)                               | Specify the single column to get max the value                                                                   |
| SelectMin     | SelectMin(col string)                               | Specify the single column to get min the value                                                                   |
| SelectCount   | SelectCount(col string)                             | Specify the single column to get total data that column                                                          |
| Where         | Where(col string, value interface{}, opt ...string) | set the condition, ex: Where("id", 1, ">") -> **"WHERE id > 1"** Where("name", "test") => **"WHERE name = 'test'"** |
| WhereIn       |                                                     |
| WhereNotIn    |                                                     |
| WhereLike     |                                                     |
| Or            |                                                     |
| OrIn          |                                                     |
| OrNotIn       |                                                     |
| OrLike        |                                                     |
| GroupBy       |                                                     |
| Join          |                                                     |
| Limit         |                                                     |
| OrderBy       |                                                     |
| Asc           |                                                     |
| Desc          |                                                     |

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
        // db tag, filled based on column name
        Id int `db:"id"`
        Name string 
        Address string
        // IsActive is boolean field
        IsActive bool `db:"is_active"`
        // BirthDate is "Date" data type column field
        BirthDate string `db:"birth_date"`
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
    if err := db.Select("name, address, birth_date").Where("is_active", true).Get(&studentList); err != nil {
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


