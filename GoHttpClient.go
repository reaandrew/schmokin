package main

import (
	"fmt"
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
	fmt.Println("Go HTTP Client Args", args)
	c := http.Client{}
	req, err := NewRequestAdapter().CreateRequest(args)
	fmt.Println("Request Method", req.Method)
	if err != nil {
		panic(err)
	}

	trace := &httptrace.ClientTrace{
		GotConn: func(connInfo httptrace.GotConnInfo) {
			fmt.Println("Got Conn")
		},
		ConnectStart: func(network, addr string) {
			fmt.Println("Dial start")
		},
		ConnectDone: func(network, addr string, err error) {
			fmt.Println("Dial done")
		},
		GotFirstResponseByte: func() {
			fmt.Println("First response byte!")
		},
		WroteHeaders: func() {
			fmt.Println("Wrote headers")
		},
		WroteRequest: func(wr httptrace.WroteRequestInfo) {
			fmt.Println("Wrote request", wr)
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
