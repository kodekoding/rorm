package rorm

import "database/sql"

// RormEngine - Engine of Raw Query ORM Library
type RormEngine interface {
	SetTableOptions(tbCaseFormat, tbPrefix string)

	SetDB(db *sql.DB)
	GetDB() *sql.DB
	GetPreparedValues() []interface{}

	GenerateRawCUDQuery(command string, data interface{})

	GetLastQuery() string
	// GetResults() []map[string]string
	// GetSingleResult() map[string]string

	Select(col ...string) *Engine
	SelectSum(col string, colAlias ...string) *Engine
	SelectAverage(col string, colAlias ...string) *Engine
	SelectMax(col string, colAlias ...string) *Engine
	SelectMin(col string, colAlias ...string) *Engine
	SelectCount(col string, colAlias ...string) *Engine

	Where(col string, value interface{}, opt ...string) *Engine
	WhereRaw(args string, value ...interface{}) *Engine

	WhereIn(col string, listOfValues ...interface{}) *Engine
	WhereNotIn(col string, listOfValues ...interface{}) *Engine
	WhereLike(col, value string) *Engine
	WhereBetween(col string, val1, val2 interface{})
	WhereNotBetween(col string, val1, val2 interface{})

	Or(col string, value interface{}, opt ...string) *Engine
	OrIn(col string, listOfValues ...interface{}) *Engine
	OrNotIn(col string, listOfValues ...interface{}) *Engine
	OrLike(col, value string) *Engine
	OrBetween(col string, val1, val2 interface{})
	OrNotBetween(col string, val1, val2 interface{})

	OrderBy(col, value string) *Engine
	Asc(col string) *Engine
	Desc(col string) *Engine

	Limit(limit int, offset ...int) *Engine
	From(tableName string) *Engine

	SQLRaw(rawQuery string, values ...interface{}) *Engine
	Get(pointerStruct interface{}) error

	Insert(data interface{}) error
	Update(data interface{}) error
	Delete(data interface{}) error
}
