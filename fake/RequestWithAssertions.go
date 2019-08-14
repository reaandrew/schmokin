package fake

import (
	"fmt"
	"log"

	schmokin "github.com/reaandrew/schmokin/core"
)

type RequestWithAssertions struct {
	request schmokin.Request
}

func (self RequestWithAssertions) IsOfType(requestType string) bool {
	result := self.request.RequestObject.Type == requestType
	if !result {
		log.Println(fmt.Sprintf("type = %s", self.request.RequestObject.Type))
	}
	return result
}
