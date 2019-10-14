package rorm

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/radityaapratamaa/rorm/lib"
)

func (re *RormEngine) generateRawCUDQuery(command string, data interface{}) {
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
		if strings.Contains(tagField.Get("rorm"), "autoincrement") {
			continue
		}
		cols += tagField.Get("rorm") + ","
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

func (re *RormEngine) Insert(data interface{}) error {
	if data == nil {
		return errors.New("Need Parameter to be passed")
	}
	re.generateRawCUDQuery("INSERT", data)
	fmt.Println(re.rawQuery)
	return nil
}

func (re *RormEngine) Update(data interface{}) error {
	if data == nil {
		return errors.New("Need Parameter to be passed")
	}
	re.generateRawCUDQuery("UPDATE", data)
	fmt.Println(re.rawQuery)
	return nil
}
