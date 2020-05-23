package account

import (
	"encoding/json"
	"io/ioutil"
)

type List struct {
	Accounts map[string]string `json:"accounts"`
	Repos map[string]string `json:"repos"`
}


func LoadAccountFromFile(filename string) (List, error) {

	if filename == "" {
		return List{}, nil
	}

	var list List

	jsonString, err := ioutil.ReadFile(filename)
	if err != nil {
		return List{}, err
	}

	err = json.Unmarshal(jsonString, &list)

	if err != nil {
		return List{}, err
	}

	return list, nil
}
