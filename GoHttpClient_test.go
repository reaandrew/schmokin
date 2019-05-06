package main_test

import (
	"testing"

	schmokin "github.com/reaandrew/schmokin"
	"github.com/stretchr/testify/assert"
)

func Test_GoHttpClient(t *testing.T) {
	client := schmokin.CreateGoHttpClient()

	response, err := client.Execute([]string{
		"-X",
		"POST",
	})

	assert.Equal(t, "POST", response.GetMethod())
}
