package main

import (
	"context"
	"log"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Schmokin(t *testing.T) {

	m := http.NewServeMux()
	s := http.Server{Addr: ":40000", Handler: m}
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
		s.Shutdown(ctx)
	}()
	m.HandleFunc("/pretty", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	t.Run("Test Status Equals", func(t *testing.T) {
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
	})

	t.Run("Test Status NotEquals", func(t *testing.T) {
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
	})

	t.Run("Test Status GreaterThan", func(t *testing.T) {
		var httpClient = CreateCurlHttpClient()
		var app = CreateSchmokinApp(httpClient)
		var args = []string{
			"http://localhost:40000/pretty",
			"--status",
			"--gt",
			"100",
		}

		var result = app.schmoke(args)
		assert.True(t, result.success)
	})

	t.Run("Test Status GreaterThanOrEqual Equal", func(t *testing.T) {
		var httpClient = CreateCurlHttpClient()
		var app = CreateSchmokinApp(httpClient)
		var args = []string{
			"http://localhost:40000/pretty",
			"--status",
			"--gte",
			"200",
		}

		var result = app.schmoke(args)
		assert.True(t, result.success)
	})

	t.Run("Test Status GreaterThanOrEqual Greater", func(t *testing.T) {
		var httpClient = CreateCurlHttpClient()
		var app = CreateSchmokinApp(httpClient)
		var args = []string{
			"http://localhost:40000/pretty",
			"--status",
			"--gte",
			"100",
		}

		var result = app.schmoke(args)
		assert.True(t, result.success)
	})

	t.Run("Test Status LessThan", func(t *testing.T) {
		var httpClient = CreateCurlHttpClient()
		var app = CreateSchmokinApp(httpClient)
		var args = []string{
			"http://localhost:40000/pretty",
			"--status",
			"--lt",
			"201",
		}

		var result = app.schmoke(args)
		assert.True(t, result.success)
	})

	t.Run("Test Status LessThanOrEqual Less", func(t *testing.T) {
		var httpClient = CreateCurlHttpClient()
		var app = CreateSchmokinApp(httpClient)
		var args = []string{
			"http://localhost:40000/pretty",
			"--status",
			"--lte",
			"201",
		}

		var result = app.schmoke(args)
		assert.True(t, result.success)
	})

	t.Run("Test Status LessThanOrEqual Equal", func(t *testing.T) {
		var httpClient = CreateCurlHttpClient()
		var app = CreateSchmokinApp(httpClient)
		var args = []string{
			"http://localhost:40000/pretty",
			"--status",
			"--lte",
			"200",
		}

		var result = app.schmoke(args)
		assert.True(t, result.success)
	})

	t.Run("Test Body Contains", func(t *testing.T) {
		var httpClient = CreateCurlHttpClient()
		var app = CreateSchmokinApp(httpClient)
		var args = []string{
			"http://localhost:40000/pretty",
			"--co",
			"OK",
		}

		var result = app.schmoke(args)
		assert.True(t, result.success)
	})
}
