package rorm

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"time"

	"github.com/radityaapratamaa/rorm/lib"
)

func (re *RormEngine) GenerateRawCUDQuery(command string, data interface{}) {
	refValue := reflect.ValueOf(data)
	tableName := ""
	re.rawQuery = command
	// cols := ""
	// values := ""
	switch refValue.Kind() {
	case reflect.Struct:
		tableName = refValue.Type().Name()
	case reflect.Ptr:
		tableName = refValue.Type().Elem().Name()
		refValue = refValue.Elem()
	}

	// Change Table Name Camel Case to Snake Case
	tableName = lib.CamelToSnakeCase(tableName)
	if command == "INSERT" {
		re.rawQuery += " INTO "
	} else if command == "DELETE" {
		re.rawQuery += " FROM "
	}
	re.rawQuery += " " + tableName + " "
	re.conditionValue = nil
	cols := "("
	values := cols
	if command == "UPDATE" {
		values = ""
	}
	for i := 0; i < refValue.NumField(); i++ {
		tagField := refValue.Type().Field(i).Tag
		if strings.Contains(tagField.Get("json"), "autoincrement") {
			continue
		}
		cols += tagField.Get("json") + ","
		if command == "INSERT" {
			values += "?,"
		} else if command == "UPDATE" {
			values += tagField.Get("rorm") + " = ?,"
		}
		re.conditionValue = append(re.conditionValue, refValue.Field(i).Interface())
	}
	cols = cols[:len(cols)-1]
	values = values[:len(values)-1]
	cols += ")"
	if command == "INSERT" {
		values += ")"
		re.rawQuery += cols + " VALUES " + values
	} else if command == "UPDATE" {
		re.rawQuery += " SET " + values
	}
	re.rawQuery = re.adjustPreparedParam(re.rawQuery)
	if re.condition != "" {
		re.convertToPreparedCondition()
		re.rawQuery += " WHERE " + re.condition
	}
}

func (re *RormEngine) executeCUDQuery(cmd string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	prepared, err := re.db.PrepareContext(ctx, re.rawQuery)
	if err != nil {
		return 0, errors.New("Error When Prepare Statement: " + err.Error())
	}
	defer prepared.Close()

	exec, err := prepared.ExecContext(ctx, re.conditionValue...)
	if err != nil {
		return 0, errors.New("Error When Execute Prepare Statement: " + err.Error())
	}

	if cmd == "INSERT" {
		return exec.LastInsertId()
	}

	return exec.RowsAffected()
}

func (re *RormEngine) Insert(data interface{}) error {
	if data == nil {
		return errors.New("Need Parameter to be passed")
	}
	command := "INSERT"
	re.GenerateRawCUDQuery(command, data)
	_, err := re.executeCUDQuery(command)

	return err
}

func (re *RormEngine) Update(data interface{}) error {
	if data == nil {
		return errors.New("Need Parameter to be passed")
	}
	command := "INSERT"
	re.GenerateRawCUDQuery(command, data)
	_, err := re.executeCUDQuery(command)
	return err
}
