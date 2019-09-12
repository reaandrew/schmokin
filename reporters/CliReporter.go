package reporters

import (
	"os"
	"text/template"

	schmokin "github.com/reaandrew/schmokin/core"
)

type CliReporter struct {
}

func (self CliReporter) Execute(result schmokin.Result) {
	reportTemplate := `
Request
----------------------------------

Headers

{{range $key, $value := .Headers}}
{{$key}} = {{$value}}
{{end}}
	`
	tmpl, err := template.New("cli").Parse(reportTemplate)

	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(os.Stdout, result)

}
