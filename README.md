# RORM (Raw Query ORM)
Raw Query ORM Library for golang (postgres, mysql)
other RDBMS coming soon

NoSQL query will be coming soon too

## Installation
```go
go get github.com/radityaapratamaa/rorm
```

## Features
| Read          | CUD    |
| :------------ | :----- |
| Select        | Insert |
| SelectSum     | Update |
| SelectAverage |        |
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

