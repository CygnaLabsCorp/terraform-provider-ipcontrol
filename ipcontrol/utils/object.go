package utils

import "strings"

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
