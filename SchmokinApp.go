package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

type SchmokinApp struct {
	httpClient SchmokinHttpClient
	targetKey  string
	target     string
	current    int
	results    ResultCollection
	data       map[string]interface{}
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

type Triple struct {
	Subject string
	Verb    string
	Object  string
}

func (instance *SchmokinApp) assertions(arg string, expected string, schmokinResult SchmokinResponse) (result Result) {
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
	case "--assert-header":
		fields := strings.Fields(expected)
		triple := Triple{
			Subject: fields[0],
			Verb:    fields[1],
			Object:  fields[2],
		}
		if strings.ToUpper(triple.Verb) == "EQ" {
			headerValue := schmokinResult.responseObj.Header.Get(triple.Subject)
			result = Result{
				Success:   headerValue == triple.Object,
				Statement: fmt.Sprintf("%s to equal %s", fmt.Sprintf("Response Header: %s", triple.Subject), triple.Object),
				Actual:    headerValue,
			}
		}
	case "--assert-context":
		fields := strings.Fields(expected)
		triple := Triple{
			Subject: fields[0],
			Verb:    fields[1],
			Object:  fields[2],
		}
		if strings.ToUpper(triple.Verb) == "EQ" {
			value := instance.data[triple.Subject]
			result = Result{
				Success:   value == triple.Object,
				Statement: fmt.Sprintf("%s to equal %s", fmt.Sprintf("Context Variable: %s", triple.Subject), triple.Object),
				Actual:    value,
			}
		}
	case "--assert-status":
		fields := strings.Fields(expected)
		if strings.ToUpper(fields[0]) == "EQ" {
			value := schmokinResult.responseObj.StatusCode
			actual, _ := strconv.Atoi(fields[1])
			result = Result{
				Success:   value == actual,
				Statement: fmt.Sprintf("%v to equal %v", fmt.Sprintf("Status Code: %v", value), actual),
				Actual:    value,
			}
		}

	}
	instance.current += 1
	return
}

func (instance *SchmokinApp) extractors(args []string, result SchmokinResponse) bool {

	switch args[instance.current] {
	case "--status":
		instance.targetKey = "HTTP Status Code"
		instance.target = strconv.Itoa(result.responseObj.StatusCode)
	case "--res-header":
		instance.checkArgs(args, instance.current, "--res-header")
		headerKey := args[instance.current+1]
		instance.target = result.responseObj.Header.Get(headerKey)
		instance.targetKey = fmt.Sprintf("Response Header: %s", headerKey)
		instance.current += 1
	case "--res-body":
		instance.target = result.payload
		instance.targetKey = "Response Body"
	case "--extract-json":
		fields := strings.Fields(args[instance.current+1])
		key := fields[0]
		json := fields[1]
		value := gjson.Get(result.payload, json)
		instance.data[key] = value.String()
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
		case "--gt", "--gte", "--lt", "--lte", "--eq", "--ne", "--co", "--assert-header", "--assert-context", "--assert-status":
			instance.checkArgs(args, instance.current, args[instance.current])
			result := instance.assertions(args[instance.current], args[instance.current+1], response)
			result.Method = response.GetMethod()
			result.Url = response.GetUrl()
			instance.addResults(result)
		case "--status", "--res-header", "--res-body", "--extract-json":
			instance.extractors(args, response)
		case "--export":
			//Need to read state file on start
			// What will the state file be called?
			//Store up the new values
			//Persist the state on completion
			state[args[instance.current+1]] = instance.target
		default:
			/*
				if instance.current > 0 && instance.current != len(args)-1 {
					panic(fmt.Sprintf("Unknown Arg: %v", args[instance.current]))
				}
			*/
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
	instance.current = 0
	var response SchmokinResponse

	var service = StateService{}
	var state = service.Load()
	var argInterceptor = CreateArgsInterceptor(state)
	args = argInterceptor.Intercept(args)

	if !strings.HasPrefix(args[0], "-") {
		response, err := instance.httpClient.Execute(args)
		log.Debug("Executing the http client")
		if err == nil {
			log.Debug("Executed the client successfully")
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
		data:       map[string]interface{}{},
	}
}
