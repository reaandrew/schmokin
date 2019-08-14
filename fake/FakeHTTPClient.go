package fake

import (
	schmokin "github.com/reaandrew/schmokin/core"
)

type FakeHTTPClient struct {
	lastRequest schmokin.Request
}

func (self *FakeHTTPClient) Execute(request schmokin.Request) schmokin.Result {
	self.lastRequest = request
	return schmokin.Result{}
}

func (self *FakeHTTPClient) Request() schmokin.Request {
	return self.lastRequest
}

func CreateFakeHTTPClient() *FakeHTTPClient {
	return &FakeHTTPClient{}
}
