package reporters

import (
	"fmt"
	"os"
	"strings"
	"text/template"

	schmokin "github.com/reaandrew/schmokin/core"
)

type CliReportViewModel struct {
	Headers map[string]string
	Timings map[string]string
}

func mapViewModel(result schmokin.Result) CliReportViewModel {
	viewModel := CliReportViewModel{
		Headers: map[string]string{},
		Timings: map[string]string{},
	}

	var paddingLength int

	for key, _ := range result.Headers {
		if len(key) > paddingLength {
			paddingLength = len(key)
		}
	}
	if paddingLength < 20 {
		paddingLength = 20
	}

	format := fmt.Sprintf("%%-%dv", paddingLength)

	for key, value := range result.Headers {
		paddedKey := fmt.Sprintf(format, key)
		viewModel.Headers[paddedKey] = strings.Join(value, ",")
	}

	viewModel.Timings[fmt.Sprintf(format, "Connect")] = result.Timings.ConnectDone.String()
	viewModel.Timings[fmt.Sprintf(format, "DNS Lookup")] = result.Timings.DnsDone.String()
	viewModel.Timings[fmt.Sprintf(format, "TLS Handshake")] = result.Timings.TlsHandshakeDone.String()
	viewModel.Timings[fmt.Sprintf(format, "Time to first byte")] = result.Timings.FirstByteDone.String()
	viewModel.Timings[fmt.Sprintf(format, "Time to complete")] = result.Timings.Complete.String()

	return viewModel
}

type CliReporter struct {
}

func (self CliReporter) Execute(result schmokin.Result) {
	reportTemplate := `
Reponse Headers
----------------------------------
{{range $key, $value := .Headers}}{{$key}} = {{$value}}
{{end}}

Timings
----------------------------------
{{range $key, $value := .Timings}}{{$key}} = {{$value}}
{{end}}`
	tmpl, err := template.New("cli").Parse(reportTemplate)

	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(os.Stdout, mapViewModel(result))
}
