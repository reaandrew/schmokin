package main_test

import (
	"testing"

	schmokin "github.com/reaandrew/schmokin"
	"github.com/stretchr/testify/assert"
)

func TestStateService(t *testing.T) {

	t.Run("Load", func(t *testing.T) {
		expected := schmokin.State{
			"body": "fubar",
		}

		schmokin.WriteGob(schmokin.StatePath, expected)

		service := schmokin.StateService{}
		state := service.Load()
		assert.Equal(t, "fubar", state["body"])
	})

	t.Run("Save", func(t *testing.T) {
		stateValue := schmokin.State{
			"body": "fubar",
		}
		service := schmokin.StateService{}
		service.Save(stateValue)

		expected := new(schmokin.State)
		schmokin.ReadGob(schmokin.StatePath, expected)
		assert.Equal(t, "fubar", (*expected)["body"])
	})
}
