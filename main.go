package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

const (
	ExpectedNotInteger string = "Argument must be a integer for the expected"
	ActualNotInteger   string = "Argument must be a integer for the actual"
	StatePath          string = "schmokin.state"
)

var (
	GreenText = color.New(color.FgGreen, color.Bold).SprintFunc()
	RedText   = color.New(color.FgRed, color.Bold).SprintFunc()
)

var SchmokinFormat = `content_type: %{content_type}\n filename_effective: %{filename_effective}\n ftp_entry_path: %{ftp_entry_path}\n http_code: %{http_code}\n http_connect: %{http_connect}\n local_ip: %{local_ip}\n local_port: %{local_port}\n num_connects: %{num_connects}\n num_redirects: %{num_redirects}\n redirect_url: %{redirect_url}\n remote_ip: %{remote_ip}\n remote_port: %{remote_port}\n size_download: %{size_download}\n size_header: %{size_header}\n size_request: %{size_request}\n size_upload: %{size_upload}\n speed_download: %{speed_download}\n speed_upload: %{speed_upload}\n ssl_verify_result: %{ssl_verify_result}\n time_appconnect: %{time_appconnect}\n time_connect: %{time_connect}\n time_namelookup: %{time_namelookup}\n time_pretransfer: %{time_pretransfer}\n time_redirect: %{time_redirect}\n time_starttransfer: %{time_starttransfer}\n time_total: %{time_total}\n url_effective: %{url_effective}\n`

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
	Results ResultCollection
}

func (self SchmokinResult) Success() bool {
	return self.Results.Success()
}

type SchmokinHttpClient interface {
	execute(args []string) SchmokinResponse
}

type SchmokinApp struct {
	httpClient SchmokinHttpClient
	targetKey  string
	target     string
	current    int
}

func SliceIndex(slice []string, predicate func(i string) bool) int {
	for i := 0; i < len(slice); i++ {
		if predicate(slice[i]) {
			return i
		}
	}
	return -1
}

func (instance *SchmokinApp) checkArgs(args []string, current int, message string) {
	if len(args) < current+2 {
		err := fmt.Errorf(fmt.Sprintf("Must supply value to compare against for %s", message))
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func (instance *SchmokinApp) assertEquality(arg string, expected string) Result {

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

func (instance *SchmokinApp) assertNumeric(arg string, expected string) Result {

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

func (instance *SchmokinApp) assertions(arg string, expected string) (result Result) {
	switch arg {
	case "--eq", "--ne":
		result = instance.assertEquality(arg, expected)
	case "--gt", "--gte", "--lt", "--lte":
		result = instance.assertNumeric(arg, expected)
	case "--co":
		result = Result{
			Actual:    instance.target,
			Statement: fmt.Sprintf("%s to contain %v", instance.targetKey, expected),
			Success:   strings.Contains(instance.target, expected),
		}
	}
	instance.current += 1
	return
}

func (instance *SchmokinApp) extractors(args []string, result SchmokinResponse) bool {

	switch args[instance.current] {
	case "--status":
		instance.targetKey = "HTTP Status Code"
		reg, _ := regexp.Compile(`http_code:\s([\d]+)`)
		result_slice := reg.FindAllStringSubmatch(result.response, -1)
		if len(result_slice) == 1 && len(result_slice[0]) == 2 {
			instance.target = result_slice[0][1]
		}
	case "--filename_effective", "--ftp_entry_path", "--http_code", "--http_connect", "--local_ip", "--local_port", "--num_connects", "--num_redirects", "--redirect_url", "--remote_ip", "--remote_port", "--size_download", "--size_header", "--size_request", "--size_upload", "--speed_download", "--speed_upload", "--ssl_verify_result", "--time_appconnect", "--time_connect", "--time_namelookup", "--time_pretransfer", "--time_redirect", "--time_starttransfer", "--time_total", "--url_effective":
		reg, _ := regexp.Compile(fmt.Sprintf(`%s:\s([\d]+)`, args[instance.current]))
		result_slice := reg.FindAllStringSubmatch(result.response, -1)
		if len(result_slice) == 1 && len(result_slice[0]) == 2 {
			instance.target = result_slice[0][1]
		}
	case "--res-header":
		instance.checkArgs(args, instance.current, "--res-header")
		regex := fmt.Sprintf(`(?i)<\s%s:\s([^\n\r]+)`, args[instance.current+1])
		reg, _ := regexp.Compile(regex)
		result_slice := reg.FindAllStringSubmatch(result.response, -1)

		if len(result_slice) == 1 && len(result_slice[0]) == 2 {
			instance.target = result_slice[0][1]
			instance.targetKey = fmt.Sprintf("Response Header: %s", args[instance.current+1])
		}
		instance.current += 1
	case "--res-body":
		instance.target = result.payload
		instance.targetKey = "Response Body"
	}

	return true
}

func (instance *SchmokinApp) schmoke(args []string) SchmokinResult {

	//Need an argument for the url

	argsToProxy := []string{}
	results := ResultCollection{}
	instance.current = 0
	var response SchmokinResponse

	var service = StateService{}
	var state = service.Load()
	var argInterceptor = CreateArgsInterceptor(state)
	args = argInterceptor.Intercept(args)

	//Still not adding the exported variable to the state file

	extraIndex := SliceIndex(args, func(i string) bool {
		return i == "--"
	})
	if extraIndex > -1 {
		argsToProxy = append(argsToProxy, args[extraIndex+1:]...)
		args = args[:extraIndex]
	}
	if !strings.HasPrefix(args[0], "-") {
		argsToProxy = append([]string{args[0]}, argsToProxy...)
		response = instance.httpClient.execute(argsToProxy)
	}

	for instance.current < len(args) {
		switch args[instance.current] {
		case "-f":
			file, err := os.Open(args[instance.current+1])
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()
			ReadLines(file, func(line string) {
				result := instance.schmoke(strings.Fields(line))
				results = append(results, result.Results...)
			})
		case "--gt", "--gte", "--lt", "--lte", "--eq", "--ne", "--co":
			instance.checkArgs(args, instance.current, args[instance.current])
			result := instance.assertions(args[instance.current], args[instance.current+1])
			result.Method = GetMethod(response)
			result.Url = GetUrl(response)
			results = append(results, result)
		case "--status", "--filename_effective", "--ftp_entry_path", "--http_code", "--http_connect", "--local_ip", "--local_port", "--num_connects", "--num_redirects", "--redirect_url", "--remote_ip", "--remote_port", "--size_download", "--size_header", "--size_request", "--size_upload", "--speed_download", "--speed_upload", "--ssl_verify_result", "--time_appconnect", "--time_connect", "--time_namelookup", "--time_pretransfer", "--time_redirect", "--time_starttransfer", "--time_total", "--url_effective", "--res-header", "--res-body":
			instance.extractors(args, response)
		case "--export":
			//Need to read state file on start
			// What will the state file be called?
			//Store up the new values
			//Persist the state on completion
			state[args[instance.current+1]] = instance.target
		default:
			if instance.current > 0 && instance.current != len(args)-1 {
				panic(fmt.Sprintf("Unknown Arg: %v", args[instance.current]))
			}
		}

		instance.current += 1
	}

	service.Save(state)

	schmokinResult := SchmokinResult{
		Results: results,
	}

	return schmokinResult
}

func GetMethod(result SchmokinResponse) string {
	regex := `(?i)>\s([\w]+)\s([^\s]+)\sHTTP`
	reg, _ := regexp.Compile(regex)
	result_slice := reg.FindAllStringSubmatch(result.response, -1)
	if len(result_slice) == 1 && len(result_slice[0]) == 3 {
		return result_slice[0][1]
	}
	return ""
}

func GetUrl(result SchmokinResponse) string {
	regex := `(?i)url_effective\:\s(.*)`
	reg, _ := regexp.Compile(regex)
	result_slice := reg.FindAllStringSubmatch(result.response, -1)
	if len(result_slice) == 1 && len(result_slice[0]) == 2 {
		return result_slice[0][1]
	}
	return ""
}

func CreateSchmokinApp(httpClient SchmokinHttpClient) *SchmokinApp {
	return &SchmokinApp{
		httpClient: httpClient,
	}
}

func ReadLines(file *os.File, visitor func(line string)) {
	reader := bufio.NewReader(file)
	for {
		var buffer bytes.Buffer

		var l []byte
		var err error
		var isPrefix bool
		for {
			l, isPrefix, err = reader.ReadLine()
			buffer.Write(l)

			// If we've reached the end of the line, stop reading.
			if !isPrefix {
				break
			}

			// If we're just at the EOF, break
			if err != nil {
				break
			}
		}

		if err == io.EOF {
			break
		}

		line := buffer.String()
		if line != "" {
			visitor(line)
		}
	}
}

func PrintResult(result SchmokinResult) {
	for _, resultItem := range result.Results {
		fmt.Println(fmt.Sprintf("%s %s", resultItem.Method, resultItem.Url))
		fmt.Println()
		fmt.Println(resultItem)
		fmt.Println()
	}
	fmt.Println()
	if result.Results.Success() {
		fmt.Println(fmt.Sprintf("Result: %s", GreenText("SUCCESS")))
	} else {
		fmt.Println(fmt.Sprintf("Result: %s", RedText("FAILURE")))
	}
}

func Run(args []string) SchmokinResult {
	var httpClient = CreateCurlHttpClient()
	var app = CreateSchmokinApp(httpClient)
	return app.schmoke(args)
}

func main() {
	result := Run(os.Args[1:])
	PrintResult(result)
}
