package schmokin

import (
	"bufio"
	"strings"

	"gopkg.in/yaml.v2"
)

type Client struct {
	httpClient HTTPClient
}

func (self Client) Execute(reader RequestReader) Result {
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

	requestObject := Request{}
	dataToDecode := strings.Join(request, "\n")
	err := yaml.Unmarshal([]byte(dataToDecode), &requestObject)

	if err != nil {
		panic(err)
	}
	requestObject.RequestObject.Data = []byte(strings.Join(data, "\n"))

	return self.httpClient.Execute(requestObject)
}

func CreateClient(httpClient HTTPClient) Client {
	return Client{
		httpClient: httpClient,
	}
}
