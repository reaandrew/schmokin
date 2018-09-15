package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

var SchmokinFormat = `content_type: %{content_type}\n filename_effective: %{filename_effective}\n ftp_entry_path: %{ftp_entry_path}\n http_code: %{http_code}\n http_connect: %{http_connect}\n local_ip: %{local_ip}\n local_port: %{local_port}\n num_connects: %{num_connects}\n num_redirects: %{num_redirects}\n redirect_url: %{redirect_url}\n remote_ip: %{remote_ip}\n remote_port: %{remote_port}\n size_download: %{size_download}\n size_header: %{size_header}\n size_request: %{size_request}\n size_upload: %{size_upload}\n speed_download: %{speed_download}\n speed_upload: %{speed_upload}\n ssl_verify_result: %{ssl_verify_result}\n time_appconnect: %{time_appconnect}\n time_connect: %{time_connect}\n time_namelookup: %{time_namelookup}\n time_pretransfer: %{time_pretransfer}\n time_redirect: %{time_redirect}\n time_starttransfer: %{time_starttransfer}\n time_total: %{time_total}\n url_effective: %{url_effective}\n`

func run() {
	processCmd := exec.Command("curl")
	//stdout, err := processCmd.StdoutPipe()
	processCmd.Start()
}

type SchmokinResponse struct {
	response string
	payload  string
}

type SchmokinResult struct {
	success bool
}

type SchmokinHttpClient interface {
	execute(args []string) SchmokinResponse
}

type CurlHttpClient struct {
	args []string
}

func (instance CurlHttpClient) execute(args []string) SchmokinResponse {
	cmd := "curl"
	if err := exec.Command(cmd, args...).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	output, err := ioutil.ReadFile("./schmokin-output")
	if err != nil {
		panic("No output file")
	}
	payload, err := ioutil.ReadFile("./schmokin-response")
	if err != nil {
		panic("No output file")
	}

	return SchmokinResponse{
		payload:  string(payload),
		response: string(output),
	}
}

func CreateCurlHttpClient(args ...string) CurlHttpClient {
	baseArgs := []string{
		"-v",
		"-s",
		fmt.Sprintf("-w '%s'", SchmokinFormat),
		"-o ./schmokin-response",
		"> ./schmokin-output",
		"2>&1",
	}
	return CurlHttpClient{
		args: append(baseArgs, args...),
	}
}

type SchmokinApp struct {
	httpClient SchmokinHttpClient
}

func (instance SchmokinApp) schmoke(request []string) SchmokinResult {

	_ = instance.httpClient.execute(request)
	return SchmokinResult{
		success: false,
	}
}

func CreateSchmokinApp(httpClient SchmokinHttpClient) SchmokinApp {
	return SchmokinApp{
		httpClient: httpClient,
	}
}

func main() {

}
