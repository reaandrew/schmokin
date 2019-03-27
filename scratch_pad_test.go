package main_test

import (
	"encoding/gob"
	"fmt"
	"os"
	"testing"
)

type State struct {
	Name string
}

func writeGob(filePath string, object interface{}) error {
	file, err := os.Create(filePath)
	defer file.Close()
	if err == nil {
		encoder := gob.NewEncoder(file)
		encoder.Encode(object)
	}
	return err
}

func readGob(filePath string, object interface{}) error {
	file, err := os.Open(filePath)
	defer file.Close()
	if err == nil {
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(object)
	}
	return err
}

func TestSomething(t *testing.T) {
	state := State{
		Name: "Andy",
	}
	err := writeGob("./schmokin.state", state)
	if err != nil {
		fmt.Println(err)
	}

	var stateRead = new(State)
	err = readGob("./schmokin.state", stateRead)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(state.Name)
	}

}
