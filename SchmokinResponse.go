package main

import "regexp"

type SchmokinResponse struct {
	response string
	payload  string
}

func (instance SchmokinResponse) GetMethod() string {
	regex := `(?i)>\s([\w]+)\s([^\s]+)\sHTTP`
	reg, _ := regexp.Compile(regex)
	result_slice := reg.FindAllStringSubmatch(instance.response, -1)
	if len(result_slice) == 1 && len(result_slice[0]) == 3 {
		return result_slice[0][1]
	}
	return ""
}

func (instance SchmokinResponse) GetUrl() string {
	regex := `(?i)url_effective\:\s(.*)`
	reg, _ := regexp.Compile(regex)
	result_slice := reg.FindAllStringSubmatch(instance.response, -1)
	if len(result_slice) == 1 && len(result_slice[0]) == 2 {
		return result_slice[0][1]
	}
	return ""
}
