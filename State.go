package main

import "strings"

type State map[string]string

func (self State) Replace(input string) string {
	result := input

	for key, value := range self {
		result = strings.Replace(result, "$"+key, value, -1)
	}

	return result
}
