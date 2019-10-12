package rorm

import rormEngine "github.com/radityaapratamaa/rorm/engine"

type EngineInterface interface {
	Column(col ...string) *rormEngine.RormEngine
	Where(col, value string, opt ...string) *rormEngine.RormEngine
	Or(col, value string, opt ...string) *rormEngine.RormEngine
	OrderBy(col, value string) *rormEngine.RormEngine
	Limit(limit int, offset ...int) *rormEngine.RormEngine
	Get(tableName string) error
}
