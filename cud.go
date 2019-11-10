package rorm

import (
	"context"
	"errors"
	"log"
	"reflect"
	"strings"

	"github.com/radityaapratamaa/rorm/constants"
	"github.com/radityaapratamaa/rorm/lib"
)

func (re *Engine) Insert(data interface{}) error {
	defer re.clearField()
	dVal := reflect.ValueOf(data)

	if err := lib.CheckDataKind(dVal, true); err != nil {
		return err
	}
	command := "INSERT"
	if dVal.Elem().Kind() == reflect.Slice {
		re.isBulk = true
	}
	// set column and preparedValue for executing data
	if re.bulkCounter == 0 {
		re.preparedData(command, data)
		re.GenerateRawCUDQuery(command, data)
	}
	_, err := re.executeCUDQuery(command)
	return err
}

func (re *Engine) BindUpdateCol(col string, otherCols ...string) *Engine {
	re.updatedCol = make(map[string]bool)
	re.updatedCol[col] = true
	if otherCols != nil {
		for _, col := range otherCols {
			re.updatedCol[col] = true
		}
	}
	return re
}

func (re *Engine) Update(data interface{}) error {
	defer re.clearField()
	dVal := reflect.ValueOf(data)
	if err := lib.CheckDataKind(dVal, false); err != nil {
		return err
	}

	command := "UPDATE"
	re.preparedData(command, data)
	re.GenerateRawCUDQuery(command, data)
	re.multiPreparedValue = nil
	re.multiPreparedValue = append(re.multiPreparedValue, re.preparedValue)
	_, err := re.executeCUDQuery(command)
	return err
}

func (re *Engine) Delete(data interface{}) error {
	defer re.clearField()
	dVal := reflect.ValueOf(data)
	if err := lib.CheckDataKind(dVal, false); err != nil {
		return err
	}
	command := "DELETE"
	re.GenerateRawCUDQuery(command, data)
	_, err := re.executeCUDQuery(command)
	return err
}

func (re *Engine) preparedData(command string, data interface{}) {
	sdValue := re.extractTableName(data)
	cols := "("
	values := "("
	if command == "UPDATE" {
		values = ""
	}
	for x := 0; x < sdValue.NumField(); x++ {
		fieldType := sdValue.Type().Field(x)
		tagField := fieldType.Tag
		if re.updatedCol != nil {
			if _, exist := re.updatedCol[tagField.Get("db")]; !exist {
				continue
			}
		}
		field := sdValue.Field(x)
		if !re.checkStructTag(tagField, field) {
			continue
		}
		col := strings.Split(tagField.Get("db"), ",")[0]
		cols += re.syntaxQuote + col + re.syntaxQuote + ","
		if command == "INSERT" {
			values += "?,"
		} else {
			values += re.syntaxQuote + col + re.syntaxQuote + " = ?,"
		}
		re.preparedValue = append(re.preparedValue, sdValue.Field(x).Interface())
	}
	re.multiPreparedValue = append(re.multiPreparedValue, re.preparedValue)

	cols = cols[:len(cols)-1]
	values = values[:len(values)-1]
	cols += ")"
	if command == "INSERT" {
		values += ")"
		re.column = cols + " VALUES " + values
	} else if command == "UPDATE" {
		re.column = " SET " + values
	}
}

func (re *Engine) GenerateRawCUDQuery(command string, data interface{}) {
	re.rawQuery = command

	re.tableName = re.options.tbPrefix + re.tableName + re.options.tbPostfix
	// Adjustment Table Name to Case Format (if available)
	switch re.options.tbFormat {
	case "camel":
		re.tableName = lib.SnakeToCamelCase(re.tableName)
	case "snake":
		re.tableName = lib.CamelToSnakeCase(re.tableName)
	}
	if command == "INSERT" {
		re.rawQuery += " INTO "
	} else if command == "DELETE" {
		re.rawQuery += " FROM "
	}
	re.rawQuery += " " + re.tableName + " " + re.column

	// re.rawQuery = re.adjustPreparedParam(re.rawQuery)
	if re.condition != "" {
		re.convertToPreparedCondition()
		re.rawQuery += " WHERE "
		re.rawQuery += re.condition
	}
	re.rawQuery = re.db.Rebind(re.rawQuery)
}

func (re *Engine) executeCUDQuery(cmd string) (int64, error) {
	// ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	// defer cancel()
	ctx := context.Background()
	prepared, err := re.db.PreparexContext(ctx, re.rawQuery)
	if err != nil {
		return 0, errors.New(constants.ErrPrepareStatement + err.Error())
	}
	defer prepared.Close()

	var affectedRows int64
	for _, pv := range re.multiPreparedValue {
		if _, err = prepared.ExecContext(ctx, pv...); err != nil {
			log.Println(err)
			return int64(0), err
		}
		affectedRows++
	}
	return affectedRows, nil
}
