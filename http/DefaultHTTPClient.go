package http

import (
	"bytes"
	"net/http"

	schmokin "github.com/reaandrew/schmokin/core"
)

type DefaultHTTPClient struct {
}

func (self DefaultHTTPClient) Execute(request schmokin.Request) schmokin.Result {
	client := &http.Client{}

	req, err := http.NewRequest(request.RequestObject.Method,
		request.RequestObject.URL,
		bytes.NewBuffer(request.RequestObject.Data))
	if err != nil {
		panic(err)
	}

	for headerKey, headerValue := range request.RequestObject.Headers {
		req.Header.Add(headerKey, headerValue)
	}
	if _, err := client.Do(req); err != nil {
		panic(err)
	}
	return schmokin.Result{}
}

func CreateDefaultHTTPClient() schmokin.HTTPClient {
	return DefaultHTTPClient{}
}
