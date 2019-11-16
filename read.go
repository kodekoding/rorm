package rorm

import (
	"context"
	"database/sql"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/kodekoding/rorm/lib"
)

func (re *Engine) GetLastQuery() string {
	return re.rawQueryBuilder.String()
}

func (re *Engine) Select(col ...string) *Engine {
	re.writeColumn(strings.Join(col, ","))
	return re
}

func (re *Engine) aggregateFuncSelect(command, col string, colAlias ...string) {
	if re.columnBuilder.String() != "" {
		re.writeColumn(",")
	}
	re.writeColumn(command)
	re.writeColumn("(")
	re.writeColumn(col)
	re.writeColumn(")")
	if colAlias != nil {
		re.writeColumn(" AS ")
		re.writeColumn(colAlias[0])
	}
}

func (re *Engine) SelectSum(col string, colAlias ...string) *Engine {
	re.aggregateFuncSelect("SUM", col, colAlias...)
	return re
}

func (re *Engine) SelectAverage(col string, colAlias ...string) *Engine {
	re.aggregateFuncSelect("AVG", col, colAlias...)

	return re
}

func (re *Engine) SelectMax(col string, colAlias ...string) *Engine {
	re.aggregateFuncSelect("MAX", col, colAlias...)

	return re
}

func (re *Engine) SelectMin(col string, colAlias ...string) *Engine {
	re.aggregateFuncSelect("MIN", col, colAlias...)

	return re
}

func (re *Engine) SelectCount(col string, colAlias ...string) *Engine {
	re.aggregateFuncSelect("COUNT", col, colAlias...)

	return re
}

func (re *Engine) SQLRaw(rawQuery string, values ...interface{}) *Engine {
	re.isRaw = true
	re.rawQuery = rawQuery
	re.rawQuery = re.adjustPreparedParam(re.rawQuery)
	re.preparedValue = values
	return re
}

func (re *Engine) From(tableName string) *Engine {
	re.tableName = tableName
	return re
}

func (re *Engine) Join(tabel, on string) *Engine {
	re.writeJoin(" JOIN ")
	re.writeJoin(tabel)
	re.writeJoin(" ON ")
	re.writeJoin(on)
	return re
}

func (re *Engine) GroupBy(col ...string) *Engine {
	re.writeGroupBy(strings.Join(col, ","))
	return re
}

func (re *Engine) Where(col string, value interface{}, opt ...string) *Engine {
	if opt != nil {
		re.generateCondition(col, value, opt[0], true)
	} else {
		re.generateCondition(col, value, "=", true)
	}
	return re
}

func (re *Engine) WhereRaw(args string, value ...interface{}) *Engine {
	if re.conditionBuilder.String() != "" {
		re.writeCondition(" AND ")
	}
	re.writeCondition(args)
	if value != nil {
		re.preparedValue = append(re.preparedValue, value...)
	}
	return re
}

func (re *Engine) WhereIn(col string, listOfValues ...interface{}) *Engine {
	value := re.generateInValue(listOfValues...)
	re.generateCondition(col, value, "IN", true)
	return re
}

func (re *Engine) WhereNotIn(col string, listOfValues ...interface{}) *Engine {
	value := re.generateInValue(listOfValues...)
	re.generateCondition(col, value, "NOT IN", true)
	return re
}

func (re *Engine) WhereBetween(col string, val1, val2 interface{}) {
	value := re.generateBetweenValue(val1, val2)
	re.generateCondition(col, value, "BETWEEN", true)
}

func (re *Engine) WhereNotBetween(col string, val1, val2 interface{}) {
	value := re.generateBetweenValue(val1, val2)
	re.generateCondition(col, value, "NOT BETWEEN", true)
}

func (re *Engine) OrBetween(col string, val1, val2 interface{}) {
	value := re.generateBetweenValue(val1, val2)
	re.generateCondition(col, value, "BETWEEN", false)
}

func (re *Engine) OrNotBetween(col string, val1, val2 interface{}) {
	value := re.generateBetweenValue(val1, val2)
	re.generateCondition(col, value, "NOT BETWEEN", false)
}

func (re *Engine) generateBetweenValue(val1, val2 interface{}) string {
	if val1 == nil || val2 == nil {
		log.Fatalln("Values cannot be nil")
	}

	ref1 := reflect.ValueOf(val1)
	ref2 := reflect.ValueOf(val2)

	if ref1.Kind() != ref2.Kind() {
		log.Fatalln("Between Values must have same datatype")
	}
	var valBuilder strings.Builder
	valBuilder.Write([]byte("'"))
	switch ref1.Kind() {
	case reflect.Int:
		valBuilder.Write([]byte(strconv.FormatInt(ref1.Int(), 10)))
	// value += strconv.FormatInt(ref1.Int(), 10)
	case reflect.String:
		valBuilder.Write([]byte(ref1.String()))
	}
	valBuilder.Write([]byte(" AND "))
	switch ref2.Kind() {
	case reflect.Int:
		valBuilder.Write([]byte(strconv.FormatInt(ref2.Int(), 10)))
	case reflect.String:
		valBuilder.Write([]byte(ref2.String()))
	}
	valBuilder.Write([]byte("'"))

	return valBuilder.String()
}

func (re *Engine) generateInValue(listValues ...interface{}) string {
	if listValues == nil {
		log.Fatalf("Values cannot be nil")
	}
	var valBuilder strings.Builder
	valBuilder.Write([]byte("("))
	// Convert all values to '....'
	for k, val := range listValues {
		reflectValue := reflect.ValueOf(val)
		valBuilder.Write([]byte("'"))
		switch reflectValue.Kind() {
		case reflect.Int:
			valBuilder.Write([]byte(strconv.FormatInt(reflectValue.Int(), 10)))
		case reflect.String:
			valBuilder.Write([]byte(reflectValue.String()))
		}
		if k < len(listValues)-1 {
			valBuilder.Write([]byte("',"))
		}
	}
	valBuilder.Write([]byte(")"))
	return valBuilder.String()
}

func (re *Engine) OrIn(col string, listOfValues ...interface{}) *Engine {
	value := re.generateInValue(listOfValues...)
	re.generateCondition(col, value, "IN", false)
	return re
}
func (re *Engine) OrNotIn(col string, listOfValues ...interface{}) *Engine {
	value := re.generateInValue(listOfValues...)
	re.generateCondition(col, value, "NOT IN", false)
	return re
}
func (re *Engine) WhereLike(col, value string) *Engine {
	re.generateCondition(col, value, "LIKE", true)
	return re
}
func (re *Engine) OrLike(col, value string) *Engine {
	re.generateCondition(col, value, "LIKE", false)
	return re
}

func (re *Engine) generateCondition(col string, nValue interface{}, opt string, isAnd bool) {
	if re.conditionBuilder.String() != "" {
		if !isAnd {
			re.writeCondition(" OR ")
		} else {
			re.writeCondition(" AND ")
		}
	}
	re.writeCondition(col)
	re.writeCondition(" ")
	re.writeCondition(opt)
	re.writeCondition(" ")
	// fmt.Println("opt " + opt)
	iValue := reflect.ValueOf(nValue)
	value := ""
	switch iValue.Kind() {
	case reflect.Int, reflect.Int64, reflect.Int8, reflect.Int16:
		value = strconv.FormatInt(iValue.Int(), 10)
	case reflect.String:
		value = iValue.String()
	case reflect.Bool:
		nBool := iValue.Bool()
		value = "1"
		if !nBool {
			value = "0"
		}
	default:
		log.Fatalln("Value is not defined")
	}
	if !strings.Contains(opt, "IN") {
		re.writeCondition("'")
		re.writeCondition(value)
		re.writeCondition("'")
	} else {
		re.writeCondition(value)
	}
}

func (re *Engine) Or(col string, value interface{}, opt ...string) *Engine {
	if opt != nil {
		re.generateCondition(col, value, opt[0], false)
	} else {
		re.generateCondition(col, value, "=", false)
	}

	return re
}

func (re *Engine) Having() {
	// coming soon
}

func (re *Engine) OrderBy(col, value string) *Engine {
	if re.orderByBuilder.String() != "" {
		re.writeOrderBy(", ")
	}
	re.writeOrderBy(col)
	re.writeOrderBy(" ")
	re.writeOrderBy(value)
	return re
}

func (re *Engine) Asc(col string) *Engine {
	re.OrderBy(col, "ASC")
	return re
}

func (re *Engine) Desc(col string) *Engine {
	re.OrderBy(col, "DESC")
	return re
}

func (re *Engine) Limit(limit int, offset ...int) *Engine {
	if offset != nil {
		re.writeLimit(strconv.Itoa(offset[0]))
		re.writeLimit(", ")
	}
	re.writeLimit(strconv.Itoa(limit))
	return re
}

// GenerateSelectQuery - Generate Select Query
func (re *Engine) GenerateSelectQuery() {
	re.rawQueryBuilder = strings.Builder{}
	if !re.isRaw {
		//===== Generated Query Start =====
		re.writeQuery("SELECT ")
		if re.column == "" {
			re.writeQuery("*")
		} else {
			re.writeQuery(re.columnBuilder.String())
		}
		re.writeQuery(" FROM ")
		re.writeQuery(re.syntaxQuote)
		re.writeQuery(re.tableName)
		re.writeQuery(re.syntaxQuote)

		if re.conditionBuilder.String() != "" {
			// Convert the Condition Value into the prepared Statement Condition
			re.convertToPreparedCondition()
			re.writeQuery(" WHERE ")
			re.writeQuery(re.conditionBuilder.String())
		}

		if re.groupByBuilder.String() != "" {
			re.writeQuery(" GROUP BY ")
			re.writeQuery(re.groupByBuilder.String())
		}

		if re.orderByBuilder.String() != "" {
			re.writeQuery(" ORDER BY ")
			re.writeQuery(re.orderByBuilder.String())
		}

		if re.limitBuilder.String() != "" {
			re.writeQuery(" LIMIT ")
			re.writeQuery(re.limitBuilder.String())
		}
	}
	re.rawQuery = re.db.Rebind(re.rawQueryBuilder.String())
}

// Get - Execute the Raw Query and get Multi Rows Result
func (re *Engine) Get(pointerStruct interface{}) error {
	defer re.clearField()
	// var err error
	dVal := reflect.ValueOf(pointerStruct)
	if err := lib.CheckDataKind(dVal, true); err != nil {
		return err
	}
	// ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	// defer cancel()
	ctx := context.Background()
	// if re.tableName, err = lib.GetStructName(pointerStruct); err != nil {
	// 	return err
	// }
	re.extractTableName(pointerStruct)

	re.GenerateSelectQuery()
	prepareVal := re.preparedValue
	return re.ExecuteSelectQuery(ctx, pointerStruct, prepareVal...)
}

func (re *Engine) ExecuteSelectQuery(ctx context.Context, pointerStruct interface{}, args ...interface{}) error {
	if re.isMultiRows {
		return re.db.SelectContext(ctx, pointerStruct, re.rawQuery, args...)
	}
	return re.db.GetContext(ctx, pointerStruct, re.rawQuery, args...)
}

func (re *Engine) scanToStructv2(rows *sql.Rows, model interface{}) error {

	mType := reflect.TypeOf(model)
	sliceElem := reflect.MakeSlice(reflect.SliceOf(mType), 0, 0).Interface()
	// sliceElem := reflect.MakeSlice(reflect.SliceOf(mType), 0, 0) //.Interface()
	log.Printf("%#v -> %#v ", sliceElem, mType)
	// sliceElem.
	return nil
}

func (re *Engine) convertToPreparedCondition() {
	// regex := regexp.MustCompile(`(LIKE .?\W+([A-Za-z0-9]+)\W.?)|((=|<|>|>=|<=|<>|!=) .?([a-zA-Z0-9]+).?)`)

	regex := regexp.MustCompile(`'(.*?)'`)
	listOfValues := regex.FindAllString(re.condition, -1)
	// matches := regex.FindAllStringSubmatch(re.condition, -1)
	re.condition = regex.ReplaceAllString(re.condition, "?")

	// re.preparedValue = nil
	for _, val := range listOfValues {
		val = strings.Replace(val, "'", "", -1)
		re.preparedValue = append(re.preparedValue, val)
	}

}
