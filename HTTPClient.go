package main

import (
	"bytes"
	"net/http"
)

type DefaultHTTPClient struct {
}

func (self DefaultHTTPClient) Execute(request SchmokinRequest) SchmokinResult {
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
	client.Do(req)
	return SchmokinResult{}
}

func CreateDefaultHTTPClient() HTTPClient {
	return DefaultHTTPClient{}
}
