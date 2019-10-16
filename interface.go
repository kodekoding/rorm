package rorm

import "database/sql"

type Rorm interface {
	SetDB(db *sql.DB)
	GetDB() *sql.DB
	GetPreparedValues() []interface{}

	GenerateRawCUDQuery(command string, data interface{})

	GetLastQuery() string
	GetResults() []map[string]string
	GetSingleResult() map[string]string

	Select(col ...string) *Engine
	SelectSum(col string, colAlias ...string) *Engine
	SelectAverage(col string, colAlias ...string) *Engine
	SelectMax(col string, colAlias ...string) *Engine
	SelectMin(col string, colAlias ...string) *Engine
	SelectCount(col string, colAlias ...string) *Engine

	Where(col, value string, opt ...string) *Engine
	WhereIn(col string, listOfValues ...interface{}) *Engine
	WhereNotIn(col string, listOfValues ...interface{}) *Engine
	WhereLike(col, value string) *Engine

	Or(col, value string, opt ...string) *Engine
	OrIn(col string, listOfValues ...interface{}) *Engine
	OrNotIn(col string, listOfValues ...interface{}) *Engine
	OrLike(col, value string) *Engine

	OrderBy(col, value string) *Engine
	Asc(col string) *Engine
	Desc(col string) *Engine

	Limit(limit int, offset ...int) *Engine
	From(tableName string) *Engine

	Raw(rawQuery string) *Engine
	Get(pointerStruct interface{}) error

	Insert(data interface{}) error
	Update(data interface{}) error
}
