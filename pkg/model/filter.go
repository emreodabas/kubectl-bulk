package model

import (
	"fmt"
	"strings"
)

func Filterlist() ([]string, error) {

	var result = []string{"filter", "multi-select"}
	return result, nil
}

func GetFilter(filter string) (string, error) {
	filterlist, _ := Filterlist()
	var filterMessage string
	filter = strings.ToLower(filter)
	for i := 0; i < len(filterlist); i++ {
		filterMessage += filterlist[i] + " "
		if strings.ToLower(filterlist[i]) == filter {
			return filterlist[i], nil
		}
	}
	return "", fmt.Errorf(filter + " is not a valid filter. Choose wisely one of them -> " + filterMessage)
}
