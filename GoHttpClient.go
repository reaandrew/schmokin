package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptrace"

	log "github.com/sirupsen/logrus"
)

type GoHttpClient struct {
}

func CreateGoHttpClient() GoHttpClient {
	return GoHttpClient{}
}

func (instance GoHttpClient) Execute(args []string) (SchmokinResponse, error) {
	c := http.Client{}
	req, err := NewRequestAdapter().CreateRequest(args)
	if err != nil {
		panic(err)
	}

	trace := &httptrace.ClientTrace{
		GotConn: func(connInfo httptrace.GotConnInfo) {
		},
		ConnectStart: func(network, addr string) {
		},
		ConnectDone: func(network, addr string, err error) {
		},
		GotFirstResponseByte: func() {
		},
		WroteHeaders: func() {
		},
		WroteRequest: func(wr httptrace.WroteRequestInfo) {
		},
	}
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	resp, err := c.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)

	return SchmokinResponse{
		payload:     bodyString,
		responseObj: resp,
	}, nil
}
