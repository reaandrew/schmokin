package main_test

import (
	"encoding/gob"
	"os"
	"testing"

	schmokin "github.com/reaandrew/schmokin"
	"github.com/stretchr/testify/assert"
)

func Test_ReadGob(t *testing.T) {
	statePath := "./schmokin.state"
	expected := schmokin.State{
		"body": "fubar",
	}
	file, err := os.Create(statePath)
	defer file.Close()
	if err == nil {
		encoder := gob.NewEncoder(file)
		encoder.Encode(expected)
	}

	var stateRead = new(schmokin.State)
	err = schmokin.ReadGob("./schmokin.state", stateRead)

	assert.Equal(t, "fubar", (*stateRead)["body"])
}

func Test_WriteGob(t *testing.T) {
	statePath := "./schmokin.state"
	expected := schmokin.State{
		"body": "fubar",
	}
	schmokin.WriteGob(statePath, expected)
	file, err := os.Open(statePath)
	defer file.Close()

	state := new(schmokin.State)
	if err == nil {
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(state)
	}

	assert.Equal(t, "fubar", (*state)["body"])
}
