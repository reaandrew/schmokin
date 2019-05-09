package main_test

import (
	"testing"

	schmokin "github.com/reaandrew/schmokin"
	"github.com/stretchr/testify/assert"
)

func TestState(t *testing.T) {
	t.Run("Replace", func(t *testing.T) {
		data := "The $speed brown $animal"

		expected := schmokin.State{
			"speed":  "quick",
			"animal": "fox",
		}

		schmokin.WriteGob(schmokin.StatePath, expected)

		service := schmokin.StateService{}
		state := service.Load()
		assert.Equal(t, "The quick brown fox", state.Replace(data))
	})
}
