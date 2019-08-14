package http_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	schmokin "github.com/reaandrew/schmokin/core"
	"github.com/reaandrew/schmokin/fake"
	schmokinHTTP "github.com/reaandrew/schmokin/http"
	"github.com/stretchr/testify/assert"
)

var request = `---
request: 
  type: http
  method: POST
  url: $URL
  headers:
    X-SOMETHING: Boom
    Content-Type: application/json
  verify: false
  pretty: true
  before:
    - ./get-reference-data.yml
---
{
  "name":"barney",
}`

func TestDefaultHttpClient(t *testing.T) {
	var requestChan = make(chan *http.Request, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestChan <- r
		fmt.Fprintln(w, "Hello, world!")
	}))
	defer server.Close()

	httpClient := schmokinHTTP.CreateDefaultHTTPClient()
	schmokin := schmokin.CreateClient(httpClient)
	schmokinRequest := strings.Replace(request, "$URL", server.URL, -1)
	reader := fake.CreateFakeRequestReader(schmokinRequest)
	result := schmokin.Execute(reader)

	select {
	case req := <-requestChan:
		assert.Equal(t, req.Method, "POST")
		assert.Equal(t, "http://"+req.Host, server.URL)
		close(requestChan)
		break
	case <-time.After(2 * time.Second):
		assert.Fail(t, "Timed out waiting for the server to respond")
	}

	assert.NotNil(t, result)
}
