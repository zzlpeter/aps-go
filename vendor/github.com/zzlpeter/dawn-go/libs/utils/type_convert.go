package utils

import (
	"encoding/json"
	"reflect"
	"strconv"
	"strings"
)

func Map2Json(m map[string]interface{}) (string, error) {
	jsonString, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(jsonString), nil
}

func Json2Map(j string) (map[string]interface{}, error) {
	b := []byte(j)
	m := make(map[string]interface{})
	err := json.Unmarshal(b, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func Struct2Json(s interface{}) (string, error) {
	b, err := json.Marshal(s)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func Json2Struct(j string, st interface{}) (interface{}, error) {
	err := json.Unmarshal([]byte(j), &st)
	if err != nil {
		return nil, err
	}
	return st, nil
}

func Struct2Map(s interface{}) map[string]interface{} {
	o1 := reflect.TypeOf(s)
	o2 := reflect.ValueOf(s)

	var m = make(map[string]interface{})
	for i := 0; i < o1.NumField(); i++ {
		tags := o1.Field(i).Tag
		var field string
		if tags.Get("map") != "" {
			field = tags.Get("map")
		} else {
			field = o1.Field(i).Name
		}
		m[field] = o2.Field(i).Interface()
	}
	return m
}

func Map2Struct(mapper map[string]interface{}, st interface{}) error {
	// make map to json
	_json, err := Map2Json(mapper)
	if err != nil {
		return err
	}
	// make json to struct
	err = json.Unmarshal([]byte(_json), &st)
	if err != nil {
		return err
	}
	return nil
}

func String2Int(s string) (int, error) {
	rst, err := strconv.Atoi(s)
	return rst, err
}

func String2Boolean(s string) bool {
	var b bool
	s = strings.ToLower(s)
	if s == "1" || s == "true" {
		b = true
	}

	return b
}