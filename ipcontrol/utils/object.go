package utils

import (
	"encoding/json"
	"fmt"
	"strings"
)

func FormatContainer(container *string) {
	if len(*container) == 0 {
		return
	}

	if strings.HasPrefix(*container, "InControl") {
		*container = "/" + *container
	} else if !strings.HasPrefix(*container, "/InControl/") {
		*container = "/InControl/" + *container
	}
}

func StructToMap(obj interface{}) (map[string]interface{}, error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	return result, err
}

func ToStringSlice(s []interface{}) (data []string, err error) {

	for _, i := range s {
		if str, ok := i.(string); ok {
			data = append(data, str)
		} else {
			err = fmt.Errorf("Expected string, got: (%v), in : (%v)", i, s)
			return
		}
	}
	return
}
