package rorm

import (
	"context"
	"errors"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/kodekoding/rorm/constants"
	"github.com/kodekoding/rorm/lib"
)

func (re *Engine) Insert(data interface{}) error {
	defer re.clearField()

	// ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	// defer cancel()
	ctx := context.Background()
	if err := re.PrepareData(ctx, "INSERT", data); err != nil {
		return err
	}
	defer re.stmt.Close()
	if _, err := re.ExecuteCUDQuery(ctx, re.preparedValue...); err != nil {
		return err
	}
	return nil
}

func (re *Engine) Update(data interface{}) error {
	defer re.clearField()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := re.PrepareData(ctx, "UPDATE", data); err != nil {
		return err
	}
	defer re.stmt.Close()
	dt := re.preparedValue
	_, err := re.ExecuteCUDQuery(ctx, dt...)
	return err
}

func (re *Engine) Delete(data interface{}) error {
	defer re.clearField()
	dVal := reflect.ValueOf(data)
	if err := lib.CheckDataKind(dVal, false); err != nil {
		return err
	}
	command := "DELETE"
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	re.GenerateRawCUDQuery(command, data)
	dt := re.preparedValue
	_, err := re.ExecuteCUDQuery(ctx, dt...)
	return err
}

func (re *Engine) PrepareData(ctx context.Context, command string, data interface{}) error {
	var err error
	dVal := reflect.ValueOf(data)
	isInsert := false
	if command == "INSERT" {
		isInsert = true
	}
	if err := lib.CheckDataKind(dVal, isInsert); err != nil {
		return err
	}
	if dVal.Elem().Kind() == reflect.Slice {
		re.PrepareMultiInsert(ctx, data)
		return nil
	}
	re.preparedData(command, data)
	re.GenerateRawCUDQuery(command, data)
	re.stmt, err = re.db.PreparexContext(ctx, re.rawQuery)
	if err != nil {
		return errors.New(constants.ErrPrepareStatement + err.Error())
	}
	return nil
}

func (re *Engine) BindUpdateCol(col string, otherCols ...string) *Engine {
	if strings.Contains(col, ",") {
		log.Fatalln("col name parameter must not contains ','")
		return nil
	}
	re.updatedCol = make(map[string]bool)
	re.updatedCol[col] = true
	if otherCols != nil {
		for _, col := range otherCols {
			re.updatedCol[col] = true
		}
	}
	return re
}
func (re *Engine) PrepareMultiInsert(ctx context.Context, data interface{}) error {
	sdValue := reflect.ValueOf(data).Elem()
	if sdValue.Len() == 0 {
		return errors.New("Data must be filled")
	}
	firstVal := sdValue.Index(0)
	if firstVal.Kind() == reflect.Ptr {
		firstVal = firstVal.Elem()
	}

	tableName := firstVal.Type().Name()
	if re.options.tbFormat == "snake" {
		tableName = lib.CamelToSnakeCase(tableName)
	} else {
		tableName = lib.SnakeToCamelCase(tableName)
	}
	tableName = re.syntaxQuote + tableName + re.syntaxQuote
	var cols strings.Builder

	cols.Write([]byte("("))
	val := cols
	for x := 0; x < firstVal.NumField(); x++ {
		colName, valid := re.getAndValidateTag(firstVal, x)
		if !valid {
			continue
		}
		fieldValue := firstVal.Field(x)
		cols.Write([]byte(re.syntaxQuote))
		cols.Write([]byte(colName))
		cols.Write([]byte(re.syntaxQuote))
		val.Write([]byte("?"))
		if x < firstVal.NumField()-1 {
			cols.Write([]byte(","))
			val.Write([]byte(","))
		}

		re.preparedValue = append(re.preparedValue, fieldValue.Interface())
	}
	cols.Write([]byte(")"))
	val.Write([]byte(")"))

	re.writeColumn(cols.String())
	re.writeColumn(" VALUES ")
	re.writeColumn(strings.Repeat(val.String()+",", sdValue.Len()-1))
	re.writeColumn(val.String())

	re.writeQuery("INSERT INTO ")
	re.writeQuery(tableName)
	re.writeQuery(" ")
	re.writeQuery(re.columnBuilder.String())
	// re.rawQuery = "INSERT INTO " + tableName + " " + re.column
	re.rawQuery = re.db.Rebind(re.rawQueryBuilder.String())
	for x := 1; x < sdValue.Len(); x++ {
		tmpVal := sdValue.Index(x).Elem()
		for z := 0; z < tmpVal.NumField(); z++ {
			fieldType := tmpVal.Type().Field(z)
			fieldValue := tmpVal.Field(z)
			if _, valid := re.checkStructTag(fieldType.Tag, fieldValue); !valid {
				continue
			}
			re.preparedValue = append(re.preparedValue, tmpVal.Field(z).Interface())
		}
	}

	var err error
	re.stmt, err = re.db.PreparexContext(ctx, re.rawQueryBuilder.String())
	if err != nil {
		return errors.New(constants.ErrPrepareStatement + err.Error())
	}
	return nil
}

func (re *Engine) getAndValidateTag(field reflect.Value, keyIndex int) (string, bool) {
	fieldType := field.Type().Field(keyIndex)
	fieldValue := field.Field(keyIndex)
	colNameTag := ""
	var valid bool
	if colNameTag, valid = re.checkStructTag(fieldType.Tag, fieldValue); !valid {

		return "", false
	}
	if colNameTag != "" {
		colNameTag = fieldType.Name
	}
	return colNameTag, true
}

func (re *Engine) preparedData(command string, data interface{}) {
	sdValue := re.extractTableName(data)

	cols := strings.Builder{}
	cols.Write([]byte("("))
	values := cols
	if command == "UPDATE" {
		values = strings.Builder{}
	}
	var valid bool
	for x := 0; x < sdValue.NumField(); x++ {
		col := ""
		if col, valid = re.getAndValidateTag(sdValue, x); !valid {
			continue
		}
		if re.updatedCol != nil {
			if _, exist := re.updatedCol[col]; !exist {
				continue
			}
		}
		cols.Write([]byte(re.syntaxQuote))
		cols.Write([]byte(col))
		cols.Write([]byte(re.syntaxQuote))
		cols.Write([]byte(","))
		if command == "INSERT" {
			values.Write([]byte("?"))
		} else {
			values.Write([]byte(re.syntaxQuote))
			values.Write([]byte(col))
			values.Write([]byte(re.syntaxQuote))
			values.Write([]byte(" = ?"))
		}
		re.preparedValue = append(re.preparedValue, sdValue.Field(x).Interface())
		if x < sdValue.NumField()-1 {
			cols.Write([]byte(","))
			values.Write([]byte(","))
		}
	}
	re.multiPreparedValue = append(re.multiPreparedValue, re.preparedValue)
	cols.Write([]byte(")"))

	re.columnBuilder = strings.Builder{}
	if command == "INSERT" {
		values.Write([]byte(")"))
		re.writeColumn(cols.String())
		re.writeColumn(" VALUES ")
	} else if command == "UPDATE" {
		re.writeColumn(" SET ")
	}
	re.writeColumn(values.String())
}

func (re *Engine) GenerateRawCUDQuery(command string, data interface{}) {
	re.rawQueryBuilder = strings.Builder{}
	re.writeQuery(command)

	tbName := strings.Builder{}
	tbName.Write([]byte(re.options.tbPrefix))
	tbName.Write([]byte(re.tableName))
	tbName.Write([]byte(re.options.tbPostfix))
	// Adjustment Table Name to Case Format (if available)
	switch re.options.tbFormat {
	case "camel":
		re.tableName = lib.SnakeToCamelCase(re.tableName)
	case "snake":
		re.tableName = lib.CamelToSnakeCase(re.tableName)
	}
	if command == "INSERT" {
		re.writeQuery(" INTO ")
	} else if command == "DELETE" {
		re.writeQuery(" FROM ")
	}

	re.writeQuery(" ")
	re.writeQuery(tbName.String())
	re.writeQuery(" ")
	re.writeQuery(re.columnBuilder.String())

	// re.rawQuery = re.adjustPreparedParam(re.rawQuery)
	if re.condition != "" {
		re.convertToPreparedCondition()
		re.writeQuery(" WHERE ")
		re.writeQuery(re.condition)
	}
	re.rawQuery = re.db.Rebind(re.rawQueryBuilder.String())
}

func (re *Engine) ExecuteCUDQuery(ctx context.Context, preparedValue ...interface{}) (int64, error) {
	var affectedRows int64
	// for _, pv := range re.multiPreparedValue {
	if _, err := re.stmt.ExecContext(ctx, preparedValue...); err != nil {
		log.Println(err)
		return int64(0), err
	}
	affectedRows++
	return affectedRows, nil
}
