package main

import (
	"fmt"
)

type StateService struct {
}

func (self StateService) Load() State {
	var stateRead = new(State)
	err := ReadGob(StatePath, stateRead)
	if err != nil {
		fmt.Println(err)
	}
	return *stateRead
}
