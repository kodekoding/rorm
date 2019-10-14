package rorm

import "database/sql"

// import rormEngine "github.com/radityaapratamaa/rorm/engine"

type EngineInterface interface {
	SetDB(db *sql.DB)
	GetDB() *sql.DB
	GetPreparedValues() []interface{}

	GenerateRawCUDQuery(command string, data interface{})

	GetLastQuery() string
	GetResults() []map[string]string
	GetSingleResult() map[string]string

	Select(col ...string) *RormEngine
	SelectSum(col string, colAlias ...string) *RormEngine
	SelectAverage(col string, colAlias ...string) *RormEngine
	SelectMax(col string, colAlias ...string) *RormEngine
	SelectMin(col string, colAlias ...string) *RormEngine
	SelectCount(col string, colAlias ...string) *RormEngine

	Where(col, value string, opt ...string) *RormEngine
	WhereIn(col string, listOfValues ...interface{}) *RormEngine
	WhereNotIn(col string, listOfValues ...interface{}) *RormEngine
	WhereLike(col, value string) *RormEngine

	Or(col, value string, opt ...string) *RormEngine
	OrIn(col string, listOfValues ...interface{}) *RormEngine
	OrNotIn(col string, listOfValues ...interface{}) *RormEngine
	OrLike(col, value string) *RormEngine

	OrderBy(col, value string) *RormEngine
	Asc(col string) *RormEngine
	Desc(col string) *RormEngine

	Limit(limit int, offset ...int) *RormEngine
	From(tableName string) *RormEngine

	Raw(rawQuery string) *RormEngine
	Get(pointerStruct interface{}) error

	Insert(data interface{}) error
	Update(data interface{}) error
}
