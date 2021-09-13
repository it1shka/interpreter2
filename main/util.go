package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func FormatStruct(structure interface{}) string {
	s, err := json.MarshalIndent(structure, "", "\t")
	if err != nil {
		panic("Failed to format structure")
	}
	return string(s)
}

func LoadFile(path string) string {
	code, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("Failed to open file %s", path)
	}
	return string(code)
}
