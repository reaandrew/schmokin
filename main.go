package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

var (
	CommitHash string
	Version    string
	BuildTime  string
)

type SchmokinRequest struct {
	RequestObject SchmokinRequestData `yaml:"request"`
}

type SchmokinRequestData struct {
	Data    []byte
	Type    string            `yaml:"type" json:"type"`
	Method  string            `yaml:"method" json:"method"`
	URL     string            `yaml:"url" json:"url"`
	Headers map[string]string `yaml:"headers"`
	Verify  bool              `yaml:"verify"`
	Pretty  bool              `yaml:"pretty"`
	Before  []string          `yaml:"before"`
	Body    []byte
}

type SchmokinResult struct {
}

type HTTPClient interface {
	Execute(request SchmokinRequest) SchmokinResult
}

type FakeHTTPClient struct {
	lastRequest SchmokinRequest
}

func (self *FakeHTTPClient) Execute(request SchmokinRequest) SchmokinResult {
	self.lastRequest = request
	return SchmokinResult{}
}

func (self *FakeHTTPClient) Request() SchmokinRequest {
	return self.lastRequest
}

type RequestWithAssertions struct {
	request SchmokinRequest
}

func (self RequestWithAssertions) IsOfType(requestType string) bool {
	result := self.request.RequestObject.Type == requestType
	if !result {
		log.Println(fmt.Sprintf("type = %s", self.request.RequestObject.Type))
	}
	return result
}

func CreateFakeHTTPClient() *FakeHTTPClient {
	return &FakeHTTPClient{}
}

type FakeRequestReader struct {
	data string
}

func (self FakeRequestReader) Read() string {
	return self.data
}

func CreateFakeRequestReader(data string) FakeRequestReader {
	return FakeRequestReader{
		data: data,
	}
}

type SchmokinRequestReader interface {
	Read() string
}

type SchmokinClient struct {
	httpClient HTTPClient
}

func (self SchmokinClient) Execute(reader SchmokinRequestReader) SchmokinResult {
	schmokinData := reader.Read()
	stringReader := strings.NewReader(schmokinData)
	scanner := bufio.NewScanner(stringReader)
	request := []string{}
	data := []string{}
	line := 0
	setData := false
	for scanner.Scan() {
		lineContent := scanner.Text()
		if line > 0 && lineContent == "---" {
			setData = true
		}

		if !setData {
			request = append(request, lineContent)
		} else {
			if lineContent != "---" {
				data = append(data, lineContent)
			}
		}
		line++
	}

	requestObject := SchmokinRequest{}
	dataToDecode := strings.Join(request, "\n")
	err := yaml.Unmarshal([]byte(dataToDecode), &requestObject)

	if err != nil {
		panic(err)
	}
	requestObject.RequestObject.Data = []byte(strings.Join(data, "\n"))

	return self.httpClient.Execute(requestObject)
}

func CreateSchmokinClient(httpClient HTTPClient) SchmokinClient {
	return SchmokinClient{
		httpClient: httpClient,
	}
}

func main() {
	f, _ := os.Open("sample2.yml")
	scanner := bufio.NewScanner(f)
	request := []string{}
	data := []string{}
	line := 0
	setData := false
	for scanner.Scan() {
		lineContent := scanner.Text()
		if line > 0 && lineContent == "---" {
			setData = true
		}

		if !setData {
			request = append(request, lineContent)
		} else {
			if lineContent != "---" {
				data = append(data, lineContent)
			}
		}
		line++
	}

	fmt.Println(strings.Join(request, "\n"))
	fmt.Println(data)
}
