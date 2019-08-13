package main_test

import (
	"testing"

	. "github.com/reaandrew/schmokin"
	"github.com/stretchr/testify/assert"
)

func TestSimpleRequest(t *testing.T) {
	httpClient := CreateFakeHTTPClient()
	reqit := CreateSchmokinClient(httpClient)
	reader := CreateFakeRequestReader(simpleRequest)
	result := reqit.Execute(reader)

	t.Run("result is not nil", func(t *testing.T) {
		assert.NotNil(t, result)
	})

	t.Run("request type is http", func(t *testing.T) {
		assert.Equal(t, httpClient.Request().RequestObject.Type, "http")
	})

	t.Run("request method is POST", func(t *testing.T) {
		assert.Equal(t, httpClient.Request().RequestObject.Method, "POST")
	})

	t.Run("request URL is https://somewhere", func(t *testing.T) {
		assert.Equal(t, httpClient.Request().RequestObject.URL, "https://somewhere")
	})

	t.Run("request Headers has X-SOMETHING Boom", func(t *testing.T) {
		assert.Equal(t, httpClient.Request().RequestObject.Headers["X-SOMETHING"], "Boom")
	})

	t.Run("request Headers has Content-Type application/json", func(t *testing.T) {
		assert.Equal(t, httpClient.Request().RequestObject.Headers["Content-Type"], "application/json")
	})

	t.Run("request verify false", func(t *testing.T) {
		assert.False(t, httpClient.Request().RequestObject.Verify)
	})

	t.Run("request pretty true", func(t *testing.T) {
		assert.True(t, httpClient.Request().RequestObject.Pretty)
	})

	t.Run("request before contains ./get-reference-data.yml", func(t *testing.T) {
		assert.Len(t, httpClient.Request().RequestObject.Before, 1)

		assert.Contains(t, httpClient.Request().RequestObject.Before, "./get-reference-data.yml")
	})

	t.Run("request body is {'name':'barney'}", func(t *testing.T) {
		assert.Equal(t, string(httpClient.Request().RequestObject.Data), `{
  "name":"barney",
}`)

	})
}
