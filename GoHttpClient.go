package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptrace"
	"os"
)

type GoHttpClient struct {
}

func CreateGoHttpClient() GoHttpClient {
	return GoHttpClient{}
}

func (instance GoHttpClient) execute(args []string) (SchmokinResponse, error) {
	c := http.Client{}
	req, err := http.NewRequest("GET", "https://www.google.com", nil)
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
	io.Copy(os.Stdout, resp.Body)
	fmt.Println("Done!")

	return SchmokinResponse{}, nil
}
