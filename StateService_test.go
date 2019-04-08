package main_test

import (
	"testing"

	schmokin "github.com/reaandrew/schmokin"
	"github.com/stretchr/testify/assert"
)

func TestStateService(t *testing.T) {

	t.Run("Loads", func(t *testing.T) {
		expected := schmokin.State{
			"body": "fubar",
		}

		schmokin.WriteGob(schmokin.StatePath, expected)

		service := schmokin.StateService{}
		state := service.Load()
		assert.Equal(t, "fubar", state["body"])
	})

}
