package rorm

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/radityaapratamaa/rorm/lib"
)

func (re *Engine) GetResults() []map[string]string {
	return re.results
}

func (re *Engine) GetLastQuery() string {
	return re.rawQuery
}

func (re *Engine) GetSingleResult() map[string]string {
	if re.results == nil {
		return nil
	}
	return re.results[0]
}

func (re *Engine) Select(col ...string) *Engine {
	re.column += strings.Join(col, ",")
	return re
}

func (re *Engine) aggregateFuncSelect(command, col string, colAlias ...string) {
	if re.column != "" {
		re.column += ","
	}
	re.column += command + "(" + col + ")"
	if colAlias != nil {
		re.column += " AS " + colAlias[0]
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
	re.join += " JOIN " + tabel + " ON " + on
	return re
}

func (re *Engine) GroupBy(col ...string) *Engine {
	re.groupBy += strings.Join(col, ",")
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
	if re.condition != "" {
		re.condition += " AND "
	}
	re.condition += args
	re.preparedValue = append(re.preparedValue, value...)
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
	value := "'"
	switch ref1.Kind() {
	case reflect.Int:
		value += strconv.FormatInt(ref1.Int(), 10)
	case reflect.String:
		value += ref1.String()
	}
	value += "' AND '"
	switch ref2.Kind() {
	case reflect.Int:
		value += strconv.FormatInt(ref2.Int(), 10)
	case reflect.String:
		value += ref2.String()
	}
	value += "'"
	return value
}

func (re *Engine) generateInValue(listValues ...interface{}) string {
	if listValues == nil {
		log.Fatalf("Values cannot be nil")
	}
	value := "("
	// Convert all values to '....'
	for _, val := range listValues {
		reflectValue := reflect.ValueOf(val)
		value += "'"
		switch reflectValue.Kind() {
		case reflect.Int:
			value += strconv.FormatInt(reflectValue.Int(), 10)
		case reflect.String:
			value += reflectValue.String()
		}
		value += "',"
	}

	// delete last ",""
	value = value[:len(value)-1]
	value += ")"
	return value
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
	if re.condition != "" {
		if !isAnd {
			re.condition += " OR "
		} else {
			re.condition += " AND "
		}
	}
	re.condition += col + " " + opt + " "
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
		re.condition += "'" + value + "'"
	} else {
		re.condition += value
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
	if re.orderBy != "" {
		re.orderBy += ", "
	}
	re.orderBy += col + " " + value
	return re
}

func (re *Engine) Asc(col string) *Engine {
	if re.orderBy != "" {
		re.orderBy += ", "
	}
	re.orderBy += col + " ASC"
	return re
}

func (re *Engine) Desc(col string) *Engine {
	if re.orderBy != "" {
		re.orderBy += ", "
	}
	re.orderBy += col + " DESC"
	return re
}

func (re *Engine) Limit(limit int, offset ...int) *Engine {
	if offset != nil {
		re.limit = strconv.Itoa(offset[0]) + ", "
	}
	re.limit += strconv.Itoa(limit)
	return re
}

// return (*sql.Stmt, error)
func (re *Engine) generateSelectQuery() {

	if !re.isRaw {
		//===== Generated Query Start =====
		re.rawQuery = "SELECT "
		if re.column == "" {
			re.rawQuery += "*"
		} else {
			re.rawQuery += re.column
		}
		re.rawQuery += " FROM "
		re.rawQuery += re.tableName

		if re.condition != "" {
			// Convert the Condition Value into the prepared Statement Condition
			re.convertToPreparedCondition()
			re.rawQuery += " WHERE "
			re.rawQuery += re.condition
		}

		if re.groupBy != "" {
			re.rawQuery += " GROUP BY "
			re.rawQuery += re.groupBy
		}

		if re.orderBy != "" {
			re.rawQuery += " ORDER BY "
			re.rawQuery += re.orderBy
		}

		if re.limit != "" {
			re.rawQuery += " LIMIT "
			re.rawQuery += re.limit
		}
	}

	// Set Prepared Raw Query
	// return re.db.Prepare(re.rawQuery)
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

	re.generateSelectQuery()

	if re.isMultiRows {
		return re.db.SelectContext(ctx, pointerStruct, re.rawQuery, re.preparedValue...)
	}
	return re.db.GetContext(ctx, pointerStruct, re.rawQuery, re.preparedValue...)

}

func (re *Engine) scanToStructv2(rows *sql.Rows, model interface{}) error {

	mType := reflect.TypeOf(model)
	sliceElem := reflect.MakeSlice(reflect.SliceOf(mType), 0, 0).Interface()
	// sliceElem := reflect.MakeSlice(reflect.SliceOf(mType), 0, 0) //.Interface()
	log.Printf("%#v -> %#v ", sliceElem, mType)
	// sliceElem.
	return nil
}

func (re *Engine) scanToStruct(rows *sql.Rows, model interface{}) error {
	v := reflect.ValueOf(model)
	if v.Kind() != reflect.Ptr {
		return errors.New("must pass a pointer, not a value, to StructScan destination") // @todo add new error message
	}

	v = reflect.Indirect(v)
	// t := v.Type()

	cols, _ := rows.Columns()

	rowCount := 0
	multiRes := []map[string]interface{}{}
	var singleRes = make(map[string]interface{})
	// columns := make([]sql.RawBytes, len(cols))
	columns := make([]interface{}, len(cols))
	columnPointers := make([]interface{}, len(cols))
	for i := range columns {
		columnPointers[i] = &columns[i]
	}
	for rows.Next() {

		if err := rows.Scan(columnPointers...); err != nil {
			return err
		}

		singleRes = make(map[string]interface{})
		for i, colName := range columns {
			var value interface{}
			value = colName
			val := reflect.TypeOf(value)
			switch val.Kind() {
			case reflect.Int64, reflect.Int:
				singleRes[cols[i]] = colName.(int64)
			default:
				singleRes[cols[i]] = string(colName.([]byte))
			}
		}
		multiRes = append(multiRes, singleRes)
		rowCount++
	}

	var willBeMarshall interface{}
	willBeMarshall = multiRes
	if len(multiRes) == 1 {
		willBeMarshall = singleRes
	} else if len(multiRes) == 0 {
		model = nil
		return nil
	}
	bJson, err := json.Marshal(willBeMarshall)
	if err != nil {
		return err
	}
	return json.Unmarshal(bJson, model)
	// }
	// for i := 0; i < v.NumField(); i++ {
	// 	field := strings.Split(t.Field(i).Tag.Get("rorm"), ",")[0]

	// 	if item, ok := m[field]; ok {
	// 		if v.Field(i).CanSet() {
	// 			if item != nil {
	// 				switch v.Field(i).Kind() {
	// 				case reflect.String:
	// 					v.Field(i).SetString(string(item.([]uint8)))
	// 				case reflect.Float32, reflect.Float64:
	// 					v.Field(i).SetFloat(item.(float64))
	// 				case reflect.Int64, reflect.Int:
	// 					v.Field(i).SetInt(item.(int64))
	// 				case reflect.Ptr:
	// 					if reflect.ValueOf(item).Kind() == reflect.Bool {
	// 						itemBool := item.(bool)
	// 						v.Field(i).Set(reflect.ValueOf(&itemBool))
	// 					}
	// 				case reflect.Struct:
	// 					v.Field(i).Set(reflect.ValueOf(item))
	// 				default:
	// 					fmt.Println(t.Field(i).Name, ": ", v.Field(i).Kind(), " - > - ", reflect.ValueOf(item).Kind()) // @todo remove after test out the Get methods
	// 				}
	// 			}
	// 		}
	// 	}
	// }

	// return nil
}

//GetRows parses recordset into map
func (re *Engine) getRows(rows *sql.Rows, pointerResult interface{}) error {
	var results []map[string]interface{}
	re.results = nil

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	// Make a slice for the values
	values := make([]sql.RawBytes, len(columns))

	// rows.Scan wants '[]interface{}' as an argument, so we must copy the
	// references into such a slice
	// See http://code.google.com/p/go-wiki/wiki/InterfaceSlice for details
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	// Fetch rows
	for rows.Next() {
		// get RawBytes from data
		err = rows.Scan(scanArgs...)
		if err != nil {
			return err
		}
		// initialize the second layer
		contents := make(map[string]interface{})

		// Now do something with the data.
		// Here we just print each column as a string.
		var value string
		for i, col := range values {
			// Here we can check if the value is nil (NULL value)
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			contents[columns[i]] = value
		}
		results = append(results, contents)
	}
	if err = rows.Err(); err != nil {
		return err
	}

	bRes, err := json.Marshal(results)
	if err != nil {
		return err
	}
	return json.Unmarshal(bRes, &pointerResult)

}

func (re *Engine) convertToPreparedCondition() {
	// regex := regexp.MustCompile(`(LIKE .?\W+([A-Za-z0-9]+)\W.?)|((=|<|>|>=|<=|<>|!=) .?([a-zA-Z0-9]+).?)`)

	regex := regexp.MustCompile(`'(.*?)'`)
	listOfValues := regex.FindAllString(re.condition, -1)
	// matches := regex.FindAllStringSubmatch(re.condition, -1)
	re.condition = regex.ReplaceAllString(re.condition, "?")

	re.condition = re.adjustPreparedParam(re.condition)

	// re.preparedValue = nil
	for _, val := range listOfValues {
		val = strings.Replace(val, "'", "", -1)
		re.preparedValue = append(re.preparedValue, val)
	}

}
