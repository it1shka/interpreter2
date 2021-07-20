package main

import "encoding/json"

func FormatStruct(structure interface{}) string {
	s, err := json.MarshalIndent(structure, "", "\t")
	if err != nil {
		panic("Failed to format structure")
	}
	return string(s)
}
