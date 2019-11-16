**Table Of Contents**
- [RORM (Raw Query ORM) Description](#RORM-Raw-Query-ORM-Description)
- [Benchmarking vs other ORM](#Benchmarking-vs-other-ORM)
- [Support Database](#Support-Database)
- [Installation](#Installation)
- [Features (will be completed soon)](#Features-will-be-completed-soon)
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
      - [All Columns](#All-Columns)
      - [Specific Column](#Specific-Column)
    - [Delete](#Delete)
# RORM (Raw Query ORM) Description
Raw Query ORM is a Query Builder as light as raw query and as easy as ORM

# Benchmarking vs other ORM

Environment: 
  - MySQL 5.7
  - MacBook Pro (mid 2017)
  - Intel Core i7-7567U @ 3.5 GHz
  - 16GB 2133 MHz LPDDR3

source : https://github.com/kihamo/orm-benchmark

command : ``` orm-benchmark -orm=all (-multi=1 default) ```

```go
Reports: 

  2000 times - Insert
       raw:     3.00s      1497568 ns/op     552 B/op     12 allocs/op
 ---  rorm:     3.18s      1588970 ns/op     472 B/op      6 allocs/op ---
       qbs:     3.56s      1780562 ns/op    4304 B/op    104 allocs/op
       orm:     3.60s      1801961 ns/op    1427 B/op     36 allocs/op
      modl:     4.41s      2206199 ns/op    1317 B/op     28 allocs/op
      hood:     4.66s      2330843 ns/op   10738 B/op    155 allocs/op
      gorp:     4.76s      2381558 ns/op    1390 B/op     29 allocs/op
      xorm:     5.40s      2699078 ns/op    2562 B/op     66 allocs/op
      gorm:     7.97s      3985524 ns/op    7695 B/op    148 allocs/op

   500 times - MultiInsert 100 row
---   rorm:     1.51s      3018679 ns/op   41888 B/op    105 allocs/op ---
       orm:     1.52s      3037726 ns/op  103987 B/op   1529 allocs/op
       raw:     2.41s      4823872 ns/op  108566 B/op    810 allocs/op
      xorm:     2.74s      5483371 ns/op  228503 B/op   4663 allocs/op
      gorp:     Not support multi insert
      hood:     Not support multi insert
      modl:     Not support multi insert
      gorm:     Not support multi insert
       qbs:     Not support multi insert

  2000 times - Update
---   rorm:     1.40s       700864 ns/op     288 B/op      5 allocs/op ---
       raw:     1.47s       734044 ns/op     616 B/op     14 allocs/op
       orm:     1.67s       836207 ns/op    1384 B/op     37 allocs/op
      xorm:     2.81s      1405482 ns/op    2665 B/op    100 allocs/op
      modl:     3.23s      1613397 ns/op    1489 B/op     36 allocs/op
       qbs:     3.27s      1637255 ns/op    4298 B/op    104 allocs/op
      gorp:     3.54s      1770792 ns/op    1536 B/op     35 allocs/op
      hood:     4.89s      2443941 ns/op   10733 B/op    155 allocs/op
      gorm:     8.99s      4494195 ns/op   18612 B/op    383 allocs/op

  4000 times - Read (will be fixed asap)
       qbs:     2.89s       721410 ns/op    6358 B/op    176 allocs/op
       raw:     3.16s       789744 ns/op    1432 B/op     37 allocs/op
       orm:     3.48s       871127 ns/op    2610 B/op     93 allocs/op
      hood:     5.80s      1450978 ns/op    4020 B/op     48 allocs/op
---   rorm:     5.96s      1489016 ns/op    1792 B/op     40 allocs/op ---
      gorm:     6.27s      1566251 ns/op   12153 B/op    239 allocs/op
      modl:     6.36s      1591131 ns/op    1873 B/op     45 allocs/op
      xorm:     6.66s      1664184 ns/op    9354 B/op    260 allocs/op
      gorp:     6.76s      1690332 ns/op    1872 B/op     52 allocs/op

  2000 times - MultiRead limit 100 (will be fixed asap)
       raw:     2.28s      1139187 ns/op   34704 B/op   1320 allocs/op
      modl:     2.32s      1159111 ns/op   49864 B/op   1721 allocs/op
       orm:     2.38s      1191493 ns/op   85020 B/op   4283 allocs/op
      gorp:     2.49s      1243186 ns/op   63685 B/op   1909 allocs/op
---   rorm:     3.75s      1875212 ns/op   40441 B/op   1536 allocs/op ---
       qbs:     3.75s      1875799 ns/op  165634 B/op   6428 allocs/op
      hood:     4.32s      2158469 ns/op  136081 B/op   6358 allocs/op
      xorm:     5.37s      2685240 ns/op  180061 B/op   8091 allocs/op
      gorm:     5.44s      2721136 ns/op  254686 B/op   6226 allocs/op
```
# Support Database
| No   | Database   |
| :--- | :--------- |
| 1    | MySQL      |
| 2    | Postgres   |
| 3    | SQL Server |

# Installation
```go
go get github.com/kodekoding/rorm
```

import to your project

```go
import "github.com/kodekoding/rorm"
```
# Features (will be completed soon)
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
        Id int `db:"id" sql:"pk"` //identify as PK column
        Name string `db:"name"`
        Address string `db:"address"`
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
#### All Columns
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
#### Specific Column
```go
    dtStudent := Student{
        Name: "change",
        IsActive: 0,
    }

    // Add Function 'BindUpdateCol' and filled parameter with db column name lists
    affected, err := db.Where("id", 1).BindUpdateCol("name", "is_active").Update(&dtStudent)
    if err != nil {
        log.Fatalln("Error When Update")
    }

    if affected > 0 {
        log.Println("Success Update")
    }
```
```sql
    -- it will generate : (mysql)
    UPDATE Student SET name = ?, is_active = ? WHERE id = ?
    -- prepared Values :
    -- ('change', 0, 1)
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


