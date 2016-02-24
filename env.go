package env

import (
	"errors"
	"os"
	"reflect"
	"strconv"
)

// Parse takes an interface{} that must be a pointer and has a type of Struct.
// If both is true it parses the environmental variables based on the tags on
// the struct's fields and sets the value for that given field. If either of
// the initial conditions are met, an error is returned depending on condition
// error that was not met.
func Parse(val interface{}) error {
	ptrRef := reflect.ValueOf(val)
	if ptrRef.Kind() != reflect.Ptr {
		return errors.New("Expected a pointer to a Struct")
	}
	ref := ptrRef.Elem()
	if ref.Kind() != reflect.Struct {
		return errors.New("Expected a struct type to be passed")
	}
	return parseEnvVars(ref, val)
}

func parseEnvVars(ref reflect.Value, val interface{}) error {
	refType := ref.Type()
	for i := 0; i < refType.NumField(); i++ {
		value := get(refType.Field(i))
		if value == "" {
			continue
		}
		if err := set(ref.Field(i), refType.Field(i), value); err != nil {
			return err
		}
	}
	return nil
}

func get(field reflect.StructField) string {
	value := os.Getenv(field.Tag.Get("env"))
	if value != "" {
		return value
	}
	return field.Tag.Get("envDefault")
}

func set(field reflect.Value, refType reflect.StructField, value string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Bool:
		bvalue, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		field.SetBool(bvalue)
	case reflect.Int:
		intValue, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return err
		}
		field.SetInt(intValue)
	default:
		return errors.New("Type is not supported")
	}
	return nil
}
