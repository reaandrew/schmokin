package main

import (
	"fmt"
	"os"
	"os/user"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"

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
		panic(err)
	}
}

func PrintResult(result SchmokinResult) {
	for _, resultItem := range result.Results {
		fmt.Println("Something")
		fmt.Println(fmt.Sprintf("%s %s", resultItem.Method, resultItem.Url))
		fmt.Println(resultItem)
		fmt.Println()
	}
	if result.Results.Success() {
		fmt.Println(fmt.Sprintf("Result: %s", GreenText("SUCCESS")))
	} else {
		fmt.Println(fmt.Sprintf("Result: %s", RedText("FAILURE")))
	}
}

func Run(args []string) SchmokinResult {
	//var httpClient = CreateCurlHttpClient()
	var httpClient = CreateGoHttpClient()
	var app = CreateSchmokinApp(httpClient)
	return app.schmoke(args)
}

func ensureWorkingDirectory() {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(usr.HomeDir)
}

func init() {
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.WarnLevel)
}

func main() {
	/*
		ensureWorkingDirectory()
		result := Run(os.Args[1:])
		if result.Error != nil {
			fmt.Println(result.Error)
		} else {
			PrintResult(result)
		}
	*/
	app := cli.NewApp()
	app.Name = "boom"
	app.Usage = "make an explosive entrance"
	app.Action = func(c *cli.Context) error {
		fmt.Println("boom! I say!")
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
