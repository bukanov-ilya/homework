package main

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type Person struct {
	Name    string `properties:"name"`
	Address string `properties:"address,omitempty"`
	Age     int    `properties:"age"`
	Married bool   `properties:"married"`
}

func Serialize(person Person) string {
	var builder strings.Builder

	personType := reflect.TypeOf(person)
	personValue := reflect.ValueOf(person)

	for i := 0; i < personType.NumField(); i++ {
		field := personType.Field(i)
		fieldValue := personValue.Field(i)

		tag := field.Tag.Get("properties")
		if tag == "" {
			continue
		}

		tagParts := strings.Split(tag, ",")
		propertyName := tagParts[0]
		hasOmitEmpty := len(tagParts) > 1 && tagParts[1] == "omitempty"

		isEmpty := isEmptyValue(fieldValue)

		if hasOmitEmpty && isEmpty {
			continue
		}

		if builder.Len() > 0 {
			builder.WriteString("\n")
		}

		valueStr := valueToString(fieldValue)

		builder.WriteString(propertyName)
		builder.WriteString("=")
		builder.WriteString(valueStr)
	}

	return builder.String()
}

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return v.String() == ""
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Bool:
		return v.Bool() == false
	default:
		return false
	}
}

func valueToString(v reflect.Value) string {
	switch v.Kind() {
	case reflect.String:
		return v.String()
	case reflect.Int:
		return strconv.Itoa(int(v.Int()))
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'f', -1, v.Type().Bits())
	case reflect.Bool:
		return strconv.FormatBool(v.Bool())
	default:
		return fmt.Sprintf("%v", v.Interface())
	}
}

func TestSerialization(t *testing.T) {
	tests := map[string]struct {
		person Person
		result string
	}{
		"test case with empty fields": {
			result: "name=\nage=0\nmarried=false",
		},
		"test case with fields": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
			},
			result: "name=John Doe\nage=30\nmarried=true",
		},
		"test case with omitempty field": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
				Address: "Paris",
			},
			result: "name=John Doe\naddress=Paris\nage=30\nmarried=true",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := Serialize(test.person)
			assert.Equal(t, test.result, result)
		})
	}
}
