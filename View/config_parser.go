package View

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Widget struct {
	Type       string           `json:"type"`
	Properties WidgetProperties `json:"properties"`
	Children   []Widget         `json:"children"`
	Action     string           `json:"action"`
	Script     string           `json:"script"`
}
type WidgetProperties struct {
	Orientation string `json:"orientation"`
	Spacing     int    `json:"spacing"`
	Label       string `json:"label"`
	// add other common properties here
}

func GetWidgetsFromConfig(filename string) ([]Widget, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var widgets []Widget
	err = json.Unmarshal(data, &widgets)
	if err != nil {
		return nil, err
	}
	fmt.Println("initializing json")
	return widgets, nil
}
