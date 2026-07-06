package env

import (
	"errors"
	"reflect"
)

func extractConfig(config any) (reflect.Value, reflect.Type, error) {
	v := reflect.ValueOf(config)
	if v.Kind() != reflect.Pointer || v.Elem().Kind() != reflect.Struct {
		return reflect.Value{}, nil, errors.New(ConfigMustBePointer)
	}
	elem := v.Elem()
	return elem, elem.Type(), nil
}

func extractFieldEnvInfo(element reflect.Value, elementType reflect.Type, index int) (*reflect.Value, string, string, bool) {
	field := element.Field(index)
	fieldType := elementType.Field(index)

	if !field.CanSet() {
		return nil, "", "", false
	}

	envKey := fieldType.Tag.Get(keyEnv)
	defaultVal := fieldType.Tag.Get(keyDefault)

	if envKey == "" {
		return nil, "", "", false
	}

	return &field, envKey, defaultVal, true
}
