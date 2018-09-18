package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Status_Equals(t *testing.T) {
	var httpClient = CreateCurlHttpClient()
	var app = CreateSchmokinApp(httpClient)
	var args = []string{
		"http://localhost:40000/pretty",
		"--status",
		"--eq",
		"200",
	}

	var result = app.schmoke(args)
	assert.True(t, result.success)
}

func Test_Status_NotEquals(t *testing.T) {
	var httpClient = CreateCurlHttpClient()
	var app = CreateSchmokinApp(httpClient)
	var args = []string{
		"http://localhost:40000/pretty",
		"--status",
		"--ne",
		"201",
	}

	var result = app.schmoke(args)
	assert.True(t, result.success)
}
