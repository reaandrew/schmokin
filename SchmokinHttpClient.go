package main

type SchmokinHttpClient interface {
	execute(args []string) SchmokinResponse
}
