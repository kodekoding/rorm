package lib

import (
	"errors"
	"log"
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
	if reflectVal.Kind() == reflect.Ptr {
		reflectVal = reflectVal.Elem()
	}
	log.Println(reflectVal.Kind().String())
	if reflectVal.Kind() == reflect.Struct {
		// log.Println("nama structnya " + reflectVal.Type().Name())
		structName = reflectVal.Type().Name()
	} else {
		typ := reflectVal.Type().Elem()
		structName = reflect.New(typ).Elem().Type().Name()
	}
	return CamelToSnakeCase(structName), nil
}

// CamelToSnakeCase - convert Camel Case to Snake Case
func CamelToSnakeCase(str string) string {
	var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}