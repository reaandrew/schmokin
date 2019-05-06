package main

type SchmokinHttpClient interface {
	Execute(args []string) (SchmokinResponse, error)
}
