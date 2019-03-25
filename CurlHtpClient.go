package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

type CurlHttpClient struct {
	args []string
}

func (instance CurlHttpClient) execute(args []string) SchmokinResponse {
	process := "curl"

	executeArgs := append(args, instance.args...)

	var output []byte
	var err error

	if output, err = exec.Command(process, executeArgs...).CombinedOutput(); err != nil {
		exitError := err.(*exec.ExitError)
		fmt.Println(string(exitError.Stderr))
		os.Exit(1)
	}

	payloadData, _ := ioutil.ReadFile("schmokin-response")

	return SchmokinResponse{
		payload:  string(payloadData),
		response: string(output),
	}
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
