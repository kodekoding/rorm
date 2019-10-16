package rorm

import (
	"context"
	"database/sql"
	"errors"
	"reflect"
	"strings"
	"time"

	"github.com/radityaapratamaa/rorm/lib"
)

func (re *Engine) Insert(data interface{}) error {
	dVal := reflect.ValueOf(data)
	if data == nil || (dVal.Kind() != reflect.Ptr) {
		return errors.New("parameter cannot be nil and must be a pointer")
	}
	command := "INSERT"
	if dVal.Elem().Kind() == reflect.Slice {
		re.isBulk = true
	}
	// set column and preparedValue for executing data
	re.preparedData(command, data)
	re.GenerateRawCUDQuery(command, data)
	_, err := re.executeCUDQuery(command)

	return err
}

func (re *Engine) Update(data interface{}) error {
	dVal := reflect.ValueOf(data)
	if data == nil || (dVal.Kind() != reflect.Ptr) || (dVal.Kind() == reflect.Slice) {
		return errors.New("parameter cannot be nil, must be a pointer, and not slice")
	}
	command := "UPDATE"
	re.preparedData(command, data)
	re.GenerateRawCUDQuery(command, data)
	_, err := re.executeCUDQuery(command)
	return err
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
	defer re.clearField()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	prepared, err := re.db.PrepareContext(ctx, re.rawQuery)
	if err != nil {
		return 0, errors.New("Error When Prepare Statement: " + err.Error())
	}
	defer prepared.Close()

	var exec sql.Result
	execErrString := "Error When Execute Prepare Statement: "
	if re.isBulk {
		for _, preparedVal := range re.multiPreparedValue {
			if exec, err = prepared.ExecContext(ctx, preparedVal...); err != nil {
				return 0, errors.New(execErrString + err.Error())
			}
		}
	} else {
		if exec, err = prepared.ExecContext(ctx, re.preparedValue...); err != nil {
			return 0, errors.New(execErrString + err.Error())
		}
	}

	if cmd == "INSERT" {
		return exec.LastInsertId()
	}

	return exec.RowsAffected()
}

func (re *Engine) preparedData(command string, data interface{}) {
	dValue := reflect.ValueOf(data).Elem()
	sdValue := dValue
	if dValue.Kind() == reflect.Slice {
		re.multiPreparedValue = nil
		for i := 0; i < dValue.Len(); i++ {
			sdValue = dValue.Index(i)
			if i == dValue.Len()-1 {
				break
			}
			if sdValue.Kind() == reflect.Ptr {
				sdValue = sdValue.Elem()
			}
			re.preparedValue = nil
			for x := 0; x < sdValue.NumField(); x++ {
				re.preparedValue = append(re.preparedValue, sdValue.Field(x).Interface())
			}
			re.multiPreparedValue = append(re.multiPreparedValue, re.preparedValue)
		}
	}

	if sdValue.Kind() == reflect.Ptr {
		sdValue = sdValue.Elem()
	}
	re.tableName = sdValue.Type().Name()

	re.preparedValue = nil
	cols := "("
	values := "("
	if command == "UPDATE" {
		values = ""
	}
	for x := 0; x < sdValue.NumField(); x++ {
		tagField := sdValue.Type().Field(x).Tag
		col := strings.Split(tagField.Get("json"), ",")[0]
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
