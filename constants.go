package main

const (
	AppHelpText = `NAME:
	{{.Name}} - {{.Usage}}
	USAGE:
	{{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}{{if .Commands}} command [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}
	{{if len .Authors}}
	AUTHOR:
	{{range .Authors}}{{ . }}{{end}}
	{{end}}{{if .Commands}}
	GLOBAL OPTIONS:
	{{range .VisibleFlags}}{{.}}
	{{end}}{{end}}{{if .Copyright }}
	COPYRIGHT:
	{{.Copyright}}
	{{end}}{{if .Version}}
	VERSION:
	{{.Version}}
	{{end}}{{if .Compiled}}
	COMPILED:
	{{.Compiled}}
	{{end}}`
)
