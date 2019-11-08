package rorm

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"time"

	"github.com/radityaapratamaa/rorm/constants"
	"github.com/radityaapratamaa/rorm/lib"
)

func (re *Engine) Insert(data interface{}) error {
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
	re.bulkCounter++
	return err
}

func (re *Engine) Update(data interface{}) error {
	dVal := reflect.ValueOf(data)
	if err := lib.CheckDataKind(dVal, false); err != nil {
		return err
	}

	command := "UPDATE"
	re.preparedData(command, data)
	re.GenerateRawCUDQuery(command, data)
	_, err := re.executeCUDQuery(command)
	return err
}

func (re *Engine) Delete(data interface{}) error {
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
	// re.preparedValue = nil
	cols := "("
	values := "("
	if command == "UPDATE" {
		values = ""
	}
	for x := 0; x < sdValue.NumField(); x++ {
		tagField := sdValue.Type().Field(x).Tag
		col := strings.Split(tagField.Get("db"), ",")[0]
		cols += col + ","
		if command == "INSERT" {
			values += "?,"
		} else {
			values += col + " = ?,"
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

	re.rawQuery = re.adjustPreparedParam(re.rawQuery)
	if re.condition != "" {
		re.convertToPreparedCondition()
		re.rawQuery += " WHERE "
		re.rawQuery += re.condition
	}
}

func (re *Engine) executeCUDQuery(cmd string) (int64, error) {
	if !re.bulkOptimized {
		defer re.clearField()
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	prepared, err := re.db.PreparexContext(ctx, re.rawQuery)
	if err != nil {
		return 0, errors.New(constants.ErrPrepareStatement + err.Error())
	}
	defer prepared.Close()

	var affectedRows int64
	for _, pv := range re.multiPreparedValue {
		if _, err = prepared.ExecContext(ctx, pv...); err != nil {
			return int64(0), err
		}
		affectedRows++
	}
	return affectedRows, nil
}
