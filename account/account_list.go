package account

import (
	"encoding/json"
	"io/ioutil"
)

type List struct {
	Accounts Accounts `json:"accounts"`
	Repos    Repos `json:"repos"`
	Reviewers Reviewers `json:"reviewers"`
}

type Accounts map[string]string
type Repos    map[string]string
type Reviewers    map[string][]string

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
