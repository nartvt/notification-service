package utils

import (
	"errors"
	"reflect"
	"strings"
)

// convert interface to Map
func ItoDictionary(value interface{}) (map[string]interface{}, error) {
	whereType := reflect.TypeOf(value).Kind()
	switch whereType {
	case reflect.Map:
		rowMap := value.(map[string]interface{})
		return rowMap, nil
	default:
		return nil, errors.New("ItoDictionary=>_DATA_ROW_TYPE_NOT_SUPPORT_")
	}
}

func BuildActivationData(data ...string) string {
	return strings.Join(data, ":")
}

func GetActivationData(data string) []string {
	return strings.Split(data, ":")
}
