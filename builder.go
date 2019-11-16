package rorm

func (re *Engine) writeQuery(str string) {
	re.rawQueryBuilder.Write([]byte(str))
}
func (re *Engine) writeColumn(str string) {
	re.columnBuilder.Write([]byte(str))
}
func (re *Engine) writeCondition(str string) {
	re.conditionBuilder.Write([]byte(str))
}
func (re *Engine) writeJoin(str string) {
	re.joinBuilder.Write([]byte(str))
}
func (re *Engine) writeGroupBy(str string) {
	re.groupByBuilder.Write([]byte(str))
}
func (re *Engine) writeLimit(str string) {
	re.limitBuilder.Write([]byte(str))
}
func (re *Engine) writeOrderBy(str string) {
	re.orderByBuilder.Write([]byte(str))
}
