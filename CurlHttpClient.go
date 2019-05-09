package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"
)

type CurlHttpClient struct {
	args []string
}

func (instance CurlHttpClient) Execute(args []string) (SchmokinResponse, error) {
	process := "curl"

	executeArgs := append(args, instance.args...)

	var output []byte
	var err error

	if output, err = exec.Command(process, executeArgs...).CombinedOutput(); err != nil {
		return SchmokinResponse{}, err
	}

	payloadData, _ := ioutil.ReadFile("schmokin-response")

	return SchmokinResponse{
		payload:  string(payloadData),
		response: string(output),
	}, nil
}

func CreateCurlHttpClient() CurlHttpClient {
	baseArgs := []string{
		"-v",
		"-s",
		fmt.Sprintf("-w '%s'", SchmokinFormat),
		"-o",
		"schmokin-response",
	}
	return CurlHttpClient{
		args: baseArgs,
	}
}
