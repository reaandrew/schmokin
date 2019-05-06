package main_test

import (
	"testing"

	schmokin "github.com/reaandrew/schmokin"
)

func Test_GoHttpClient(t *testing.T) {

	server := CreateTestServer()
	defer server.Stop()
	go server.Start()

	t.Run("POST", func(t *testing.T) {
		client := schmokin.CreateGoHttpClient()

		client.Execute([]string{
			"http://localhost:40000/pretty",
			"-X",
			"POST",
		})

		//assert.Equal(t, "POST", response.GetMethod())
	})
}
