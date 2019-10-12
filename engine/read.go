package rorm

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func (re *RormEngine) GetResults() []map[string]string {
	return re.results
}

func (re *RormEngine) GetSingleResult() map[string]string {
	return re.result
}

func (re *RormEngine) Column(col ...string) *RormEngine {
	re.column += strings.Join(col, ",")
	return re
}

func (re *RormEngine) Where(col, value string, opt ...string) *RormEngine {

	if re.condition != "" {
		re.condition += " AND "
	}
	re.condition += col
	if opt != nil {
		re.condition += " " + opt[0]
	} else {
		re.condition += " = "
	}
	re.condition += value

	return re
}

func (re *RormEngine) Or(col, value string, opt ...string) *RormEngine {

	if re.condition != "" {
		re.condition += " OR "
	}
	re.condition += col
	if opt != nil {
		re.condition += " " + opt[0]
	} else {
		re.condition += " = "
	}
	re.condition += value

	return re
}

func (re *RormEngine) OrderBy(col, value string) *RormEngine {
	if re.orderBy != "" {
		re.orderBy += ", "
	}
	re.orderBy += col + " " + value
	return re
}

func (re *RormEngine) Limit(limit int, offset ...int) *RormEngine {
	if offset != nil {
		re.limit = strconv.Itoa(offset[0]) + ", "
	}
	re.limit += strconv.Itoa(limit)
	return re
}

func (re *RormEngine) Get(tableName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	re.rawQuery = "SELECT "
	if re.column == "" {
		re.rawQuery += "*"
	} else {
		re.rawQuery += re.column
	}
	re.rawQuery += " FROM " + tableName

	if re.condition != "" {
		re.convertToPreparedCondition()
		re.rawQuery += " WHERE " + re.condition
	}

	if re.orderBy != "" {
		re.rawQuery += " ORDER BY " + re.orderBy
	}

	if re.limit != "" {
		re.rawQuery += " LIMIT " + re.limit
	}

	prepared, err := re.DB.Prepare(re.rawQuery)
	if err != nil {
		return errors.New("Error When Prepared Query: " + err.Error())
	}
	defer prepared.Close()

	exec, err := prepared.QueryContext(ctx, re.conditionValue)
	if err != nil {
		return errors.New("Error When Execute Prepared Statement: " + err.Error())
	}
	defer exec.Close()

	err = re.getRows(exec)
	if err != nil {
		return errors.New("Error When Get Rows: " + err.Error())
	}

	return nil
}

//GetRows parses recordset into map
func (re *RormEngine) getRows(rows *sql.Rows) error {
	var results []map[string]string

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
		contents := make(map[string]string)

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
			results = append(results, contents)
		}
	}
	if err = rows.Err(); err != nil {
		return err
	}
	re.results = results
	return nil
}

func (re *RormEngine) convertToPreparedCondition() {
	tmpCond := re.condition
	regex := regexp.MustCompile(`= '(.*?)'|= .(.*?).*`)
	re.condition = regex.ReplaceAllString(re.condition, "= ?")

	reg2 := regexp.MustCompile(`[a-z]+.?= `)
	tmpCond = reg2.ReplaceAllString(tmpCond, "")
	reg3 := regexp.MustCompile(`.?(AND|OR).?`)
	tmpCond = reg3.ReplaceAllString(tmpCond, ",")
	re.conditionValue = strings.Split(tmpCond, ",")
}
