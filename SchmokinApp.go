package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

type SchmokinApp struct {
	httpClient SchmokinHttpClient
	targetKey  string
	target     string
	current    int
	results    ResultCollection
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

	log.Debug(instance.target)
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

func (instance *SchmokinApp) processArgs(args []string, response SchmokinResponse, state State) {
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
				log.WithField("result_count", len(result.Results)).Debug("File Line Executed")
				instance.addResults(result.Results...)
			})
		case "--gt", "--gte", "--lt", "--lte", "--eq", "--ne", "--co":
			instance.checkArgs(args, instance.current, args[instance.current])
			result := instance.assertions(args[instance.current], args[instance.current+1])
			result.Method = response.GetMethod()
			result.Url = response.GetUrl()
			instance.addResults(result)
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
}

func (instance *SchmokinApp) addResults(results ...Result) {
	log.Debug("Adding Result")
	instance.results = append(instance.results, results...)
}

func (instance *SchmokinApp) schmoke(args []string) SchmokinResult {
	if len(args) == 0 {
		return SchmokinResult{
			Error: fmt.Errorf("Must supply arguments"),
		}
	}
	//Need an argument for the url
	argsToProxy := []string{}
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
		response, err := instance.httpClient.Execute(argsToProxy)
		log.Debug("Executing the curl")
		if err == nil {
			log.Debug("Executed the curl successfully")
			instance.processArgs(args, response, state)
		} else {
			log.Error(err)
			instance.addResults(Result{
				Error: err,
			})
		}
	} else {
		instance.processArgs(args, response, state)
	}

	service.Save(state)

	schmokinResult := SchmokinResult{
		Results: instance.results,
	}

	return schmokinResult
}

func CreateSchmokinApp(httpClient SchmokinHttpClient) *SchmokinApp {
	return &SchmokinApp{
		httpClient: httpClient,
		results:    ResultCollection{},
	}
}
