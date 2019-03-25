package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

const (
	ExpectedNotInteger string = "Argument must be a integer for the expected"
	ActualNotInteger   string = "Argument must be a integer for the actual"
)

var (
	green = color.New(color.FgGreen).SprintFunc()
	red   = color.New(color.FgRed).SprintFunc()
)

var SchmokinFormat = `content_type: %{content_type}\n filename_effective: %{filename_effective}\n ftp_entry_path: %{ftp_entry_path}\n http_code: %{http_code}\n http_connect: %{http_connect}\n local_ip: %{local_ip}\n local_port: %{local_port}\n num_connects: %{num_connects}\n num_redirects: %{num_redirects}\n redirect_url: %{redirect_url}\n remote_ip: %{remote_ip}\n remote_port: %{remote_port}\n size_download: %{size_download}\n size_header: %{size_header}\n size_request: %{size_request}\n size_upload: %{size_upload}\n speed_download: %{speed_download}\n speed_upload: %{speed_upload}\n ssl_verify_result: %{ssl_verify_result}\n time_appconnect: %{time_appconnect}\n time_connect: %{time_connect}\n time_namelookup: %{time_namelookup}\n time_pretransfer: %{time_pretransfer}\n time_redirect: %{time_redirect}\n time_starttransfer: %{time_starttransfer}\n time_total: %{time_total}\n url_effective: %{url_effective}\n`

func run() {
	processCmd := exec.Command("curl")
	//stdout, err := processCmd.StdoutPipe()
	processCmd.Start()
}

func checkErr(err error, msg string) {
	if err != nil {
		err = fmt.Errorf(msg)
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

type SchmokinResponse struct {
	response string
	payload  string
}

type SchmokinResult struct {
	success bool
	Url     string
	Method  string
	Results ResultCollection
}

type SchmokinHttpClient interface {
	execute(args []string) SchmokinResponse
}

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

type Result struct {
	Success   bool
	Statement string
	Expected  interface{}
	Actual    interface{}
}

func (instance Result) String() string {

	if instance.Success {
		statement := fmt.Sprintf("Expect %s", instance.Statement)
		return fmt.Sprintf("%s : %s", green("PASS"), statement)
	} else {
		statement := fmt.Sprintf("Expected %s actual %s", instance.Statement, instance.Actual)
		return fmt.Sprintf("%s : %s", red("FAIL"), statement)
	}
}

type SchmokinApp struct {
	httpClient SchmokinHttpClient
	targetKey  string
	target     string
}

func SliceIndex(slice []string, predicate func(i string) bool) int {
	for i := 0; i < len(slice); i++ {
		if predicate(slice[i]) {
			return i
		}
	}
	return -1
}

func (instance SchmokinApp) checkArgs(args []string, current int, message string) {
	if len(args) < current+2 {
		err := fmt.Errorf(fmt.Sprintf("Must supply value to compare against for %s", message))
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func (instance SchmokinApp) assertEquality(arg string, expected string) Result {

	result := Result{
		Actual: instance.target,
	}

	switch arg {
	case "--eq":
		result.Statement = fmt.Sprintf("%s to equal %s", instance.targetKey, expected)
		result.Success = expected == instance.target
	case "--ne":
		result.Statement = fmt.Sprintf("%s to equal %s", instance.targetKey, expected)
		result.Success = expected != instance.target
	}

	return result
}

func (instance SchmokinApp) assertNumeric(arg string, expected string) Result {

	expectedValue, err := strconv.Atoi(expected)
	checkErr(err, ExpectedNotInteger)
	actual, err := strconv.Atoi(instance.target)
	checkErr(err, ActualNotInteger)
	result := Result{
		Actual: actual,
	}
	switch arg {
	case "--gt":
		result.Statement = fmt.Sprintf("%s is greater than %v", instance.targetKey, expected)
		result.Success = actual > expectedValue
	case "--gte":
		result.Statement = fmt.Sprintf("%s is greater than or equal %v", instance.targetKey, expected)
		result.Success = actual >= expectedValue
	case "--lt":
		result.Statement = fmt.Sprintf("%s is less than %v", instance.targetKey, expected)
		result.Success = actual < expectedValue
	case "--lte":
		result.Statement = fmt.Sprintf("%s is less than or equal %v", instance.targetKey, expected)
		result.Success = actual <= expectedValue
	}

	return result
}

func (instance SchmokinApp) schmoke(args []string) SchmokinResult {

	argsToProxy := []string{args[0]}
	extraIndex := SliceIndex(args, func(i string) bool {
		return i == "--"
	})
	if extraIndex > -1 {
		argsToProxy = append(argsToProxy, args[extraIndex+1:]...)
		args = args[:extraIndex]
	}

	result := instance.httpClient.execute(argsToProxy)

	results := ResultCollection{}

	success := true
	current := 0

	for current < len(args) {
		switch args[current] {
		case "--status":
			instance.targetKey = "HTTP Status Code"
			reg, _ := regexp.Compile(`http_code:\s([\d]+)`)
			result_slice := reg.FindAllStringSubmatch(result.response, -1)
			if len(result_slice) == 1 && len(result_slice[0]) == 2 {
				instance.target = result_slice[0][1]
			}
		case "--filename_effective", "--ftp_entry_path", "--http_code", "--http_connect", "--local_ip", "--local_port", "--num_connects", "--num_redirects", "--redirect_url", "--remote_ip", "--remote_port", "--size_download", "--size_header", "--size_request", "--size_upload", "--speed_download", "--speed_upload", "--ssl_verify_result", "--time_appconnect", "--time_connect", "--time_namelookup", "--time_pretransfer", "--time_redirect", "--time_starttransfer", "--time_total", "--url_effective":
			reg, _ := regexp.Compile(fmt.Sprintf(`%s:\s([\d]+)`, args[current]))
			result_slice := reg.FindAllStringSubmatch(result.response, -1)
			if len(result_slice) == 1 && len(result_slice[0]) == 2 {
				instance.target = result_slice[0][1]
			}
		case "--eq", "--ne":
			instance.checkArgs(args, current, args[current])
			var result = instance.assertEquality(args[current], args[current+1])
			results = append(results, result)
			current += 1
		case "--gt", "--gte", "--lt", "--lte":
			instance.checkArgs(args, current, args[current])
			var result = instance.assertNumeric(args[current], args[current+1])
			results = append(results, result)
			current += 1
		case "--co":
			//TODO: Use --co with other parameters
			instance.checkArgs(args, current, "--co")
			var expected = args[current+1]
			results = append(results, Result{
				Actual:    instance.target,
				Statement: fmt.Sprintf("%s to contain %v", instance.targetKey, expected),
				Success:   strings.Contains(instance.target, expected),
			})
			current += 1
		case "--res-header":
			instance.checkArgs(args, current, "--res-header")
			regex := fmt.Sprintf(`(?i)<\s%s:\s([^\n\r]+)`, args[current+1])
			reg, _ := regexp.Compile(regex)
			result_slice := reg.FindAllStringSubmatch(result.response, -1)

			if len(result_slice) == 1 && len(result_slice[0]) == 2 {
				instance.target = result_slice[0][1]
				instance.targetKey = fmt.Sprintf("Response Header: %s", args[current+1])
			}
			current += 1
		case "--res-body":
			instance.target = result.payload
			instance.targetKey = "Response Body"
		default:
			if current > 0 {
				panic(fmt.Sprintf("Unknown Arg: %v", args[current]))
			}
		}

		current += 1
	}

	schmokinResult := SchmokinResult{
		success: success,
		Results: results,
	}

	regex := `(?i)>\s([\w]+)\s([^\s]+)\sHTTP`
	reg, _ := regexp.Compile(regex)
	result_slice := reg.FindAllStringSubmatch(result.response, -1)
	if len(result_slice) == 1 && len(result_slice[0]) == 3 {
		schmokinResult.Method = result_slice[0][1]
	}

	regex = `(?i)url_effective\:\s(.*)`
	reg, _ = regexp.Compile(regex)
	result_slice = reg.FindAllStringSubmatch(result.response, -1)
	if len(result_slice) == 1 && len(result_slice[0]) == 2 {
		schmokinResult.Url = result_slice[0][1]
	}

	return schmokinResult
}

func CreateSchmokinApp(httpClient SchmokinHttpClient) SchmokinApp {
	return SchmokinApp{
		httpClient: httpClient,
	}
}

func main() {
	var httpClient = CreateCurlHttpClient()
	var app = CreateSchmokinApp(httpClient)
	var result = app.schmoke(os.Args[1:])

	fmt.Println(fmt.Sprintf("%s %s", result.Method, result.Url))
	fmt.Println()
	for _, resultItem := range result.Results {
		fmt.Println(resultItem)
	}
	fmt.Println()
	if result.Results.Success() {
		fmt.Println(fmt.Sprintf("Result: %s", green("SUCCESS")))
	} else {
		fmt.Println(fmt.Sprintf("Result: %s", red("FAILURE")))
	}
}
