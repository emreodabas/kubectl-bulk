package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func Actionlist() ([]Action, error) {
	// Open our jsonFile
	jsonFile, err := os.Open("pkg/model/actions.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println("ERROR!!")
		fmt.Println(err)
		return []Action{}, err
	}

	defer jsonFile.Close()
	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)
	// we initialize our Users array
	var actions Actions

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	err = json.Unmarshal(byteValue, &actions)
	if err != nil {
		fmt.Errorf("Marshall problem")
	}
	return actions.Actions, nil
}

type Actions struct {
	Actions []Action `json:"actions"`
}

type Action struct {
	Name         string `json:"name"`
	Descriptions string `json:"desc"`
}

func GetAction(action string) (Action, error) {
	actionlist, _ := Actionlist()
	var actionMessage string
	action = strings.ToLower(action)
	for i := 0; i < len(actionlist); i++ {
		actionMessage += actionlist[i].Name + " "
		if strings.ToLower(actionlist[i].Name) == action {
			return actionlist[i], nil
		}
	}
	return Action{}, fmt.Errorf(action + " is not a valid action. Choose wisely one of them -> " + actionMessage)
}
