package http

import (
	"bytes"
	"crypto/tls"
	"net/http"
	"net/http/httptrace"
	"time"

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

	var start, connect, dns, tlsHandshake time.Time
	var connectDone, dnsDone, tlsHandshakeDone, firstByteDone, doneDone time.Duration

	trace := &httptrace.ClientTrace{
		DNSStart:             func(dsi httptrace.DNSStartInfo) { dns = time.Now() },
		DNSDone:              func(ddi httptrace.DNSDoneInfo) { dnsDone = time.Since(dns) },
		TLSHandshakeStart:    func() { tlsHandshake = time.Now() },
		TLSHandshakeDone:     func(cs tls.ConnectionState, err error) { tlsHandshakeDone = time.Since(tlsHandshake) },
		ConnectStart:         func(network, addr string) { connect = time.Now() },
		ConnectDone:          func(network, addr string, err error) { connectDone = time.Since(connect) },
		GotFirstResponseByte: func() { firstByteDone = time.Since(start) },
	}

	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	start = time.Now()
	if response, err := client.Do(req); err != nil {
		panic(err)
	} else {
		doneDone = time.Since(start)
		result := schmokin.NewResult()
		for key, value := range response.Header {
			result.Headers[key] = value
		}
		result.Timings = schmokin.Timings{
			ConnectDone:      connectDone,
			DnsDone:          dnsDone,
			FirstByteDone:    firstByteDone,
			TlsHandshakeDone: tlsHandshakeDone,
			Complete:         doneDone,
		}
		return result
	}
}

func CreateDefaultHTTPClient() schmokin.HTTPClient {
	return DefaultHTTPClient{}
}
