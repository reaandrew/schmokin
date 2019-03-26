package main

import "fmt"

type Result struct {
	Url       string
	Method    string
	Success   bool
	Statement string
	Actual    interface{}
}

func (instance Result) String() string {

	if instance.Success {
		statement := fmt.Sprintf("Expect %s", instance.Statement)
		return fmt.Sprintf("%s : %s", GreenText("PASS"), statement)
	} else {
		statement := fmt.Sprintf("Expected %s actual %s", instance.Statement, instance.Actual)
		return fmt.Sprintf("%s : %s", RedText("FAIL"), statement)
	}
}
