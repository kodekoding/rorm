package lib

import (
	"errors"
	"reflect"
	"regexp"
	"strings"
)

func GetStructName(pointerStruct interface{}) (string, error) {
	structName := ""
	reflectVal := reflect.ValueOf(pointerStruct)
	if reflectVal.Kind() != reflect.Ptr {
		return "", errors.New("The Struct must be pointer")
	}
	reflectVal = reflectVal.Elem()

	if reflectVal.Kind() == reflect.Struct {
		// log.Println("nama structnya " + reflectVal.Type().Name())
		structName = reflectVal.Type().Name()
	} else {
		typ := reflectVal.Type().Elem()
		structName = reflect.New(typ).Elem().Type().Name()
	}
	return CamelToSnakeCase(structName), nil
}

func IssetSliceKey(arr interface{}, index int) bool {
	s := reflect.ValueOf(arr)
	if reflect.TypeOf(arr).Kind() != reflect.Slice {
		return false
	}
	return (s.Len() > index)
}

// CheckDataKind - Check data parameter when Create Update Delete Query
func CheckDataKind(data reflect.Value, isInsert bool) error {
	if data.Kind() != reflect.Ptr {
		return errors.New("data is not pointer")
	}
	if !isInsert && data.Elem().Kind() == reflect.Slice {
		return errors.New("data cannot be slice")
	}
	return nil
}

// CamelToSnakeCase - convert Camel Case to Snake Case
func CamelToSnakeCase(str string) string {
	var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

// SnakeToCamelCase - convert Snake Case to Camel Case
func SnakeToCamelCase(str string) string {
	var link = regexp.MustCompile("(^[A-Za-z])|_([A-Za-z])")
	return link.ReplaceAllStringFunc(str, func(s string) string {
		return strings.ToUpper(strings.Replace(s, "_", "", -1))
	})
}
