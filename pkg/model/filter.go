package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Filters struct {
	Filters []Filter `json:"actions"`
}

type Filter struct {
	Name           string `json:"name"`
	Description    string `json:"desc"`
	Sample         string `json:"sample"`
	NeedFilterable bool   `json:"needFilter"`
}

func Filterlist(verbs []string) ([]Filter, error) {

	// Open our jsonFile
	jsonFile, err := os.Open("pkg/model/filters.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println("ERROR!!")
		fmt.Println(err)
		return []Filter{}, err
	}

	defer jsonFile.Close()
	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)
	// we initialize our Users array
	var filters Filters

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	err = json.Unmarshal(byteValue, &filters)
	if err != nil {
		fmt.Errorf("Marshall problem")
	}

	return filters.Filters, nil
}

//func GetFilter(filter string) (string, error) {
//	filterlist, _ := Filterlist(nil)
//	var filterMessage string
//	filter = strings.ToLower(filter)
//	for i := 0; i < len(filterlist); i++ {
//		filterMessage += filterlist[i] + " "
//		if strings.ToLower(filterlist[i]) == filter {
//			return filterlist[i], nil
//		}
//	}
//	return "", fmt.Errorf(filter + " is not a valid filter. Choose wisely one of them -> " + filterMessage)
//}
