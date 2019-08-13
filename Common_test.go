package main_test

import "strings"

var simpleRequest = `---
request: 
  type: http
  method: POST
  url: https://somewhere
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

func CreateRequest(url string) string {
	return strings.Replace(simpleRequest, "https://somewhere", url, -1)
}
