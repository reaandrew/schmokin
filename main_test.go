package main_test

import (
	"testing"

	. "github.com/reaandrew/schmokin"
)

func TestOutputVersion(t *testing.T) {
	expectedVersion := "0.1.0"
	Version = expectedVersion
	//cmd := exec.Command("sh", "-c", "echo stdout; echo 1>&2 stderr")

	//check if the app exists first and fail the build if it does not

	Execute([]string{""})
}
