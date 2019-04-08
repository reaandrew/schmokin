package main_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	schmokin "github.com/reaandrew/schmokin"
)

func Test_Schmokin(t *testing.T) {
	m := http.NewServeMux()
	s := http.Server{Addr: ":40000", Handler: m}
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
		s.Shutdown(ctx)
	}()
	m.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			body = []byte("not set")
		}
		w.Write(body)
	})
	m.HandleFunc("/pretty", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("X-FU", "BAR")
		if r.Method == http.MethodGet {
			w.Write([]byte("OK"))
		} else {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				body = []byte("not set")
			}
			message := fmt.Sprintf("Method: %v Body: %v", r.Method, string(body))
			w.Write([]byte(message))
		}
	})
	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	t.Run("Test Status Equals", func(t *testing.T) {
		var args = []string{
			"http://localhost:40000/pretty",
			"--status",
			"--eq",
			"200",
		}

		var result = schmokin.Run(args)
		assert.True(t, result.Success())
	})

	t.Run("Test Status NotEquals", func(t *testing.T) {
		var args = []string{
			"http://localhost:40000/pretty",
			"--status",
			"--ne",
			"201",
		}

		var result = schmokin.Run(args)
		assert.True(t, result.Success())
	})

	t.Run("Test Status GreaterThan", func(t *testing.T) {
		var args = []string{
			"http://localhost:40000/pretty",
			"--status",
			"--gt",
			"100",
		}

		var result = schmokin.Run(args)
		assert.True(t, result.Success())
	})

	t.Run("Test Status GreaterThanOrEqual Equal", func(t *testing.T) {
		var args = []string{
			"http://localhost:40000/pretty",
			"--status",
			"--gte",
			"200",
		}

		var result = schmokin.Run(args)
		assert.True(t, result.Success())
	})

	t.Run("Test Status GreaterThanOrEqual Greater", func(t *testing.T) {
		var args = []string{
			"http://localhost:40000/pretty",
			"--status",
			"--gte",
			"100",
		}

		var result = schmokin.Run(args)
		assert.True(t, result.Success())
	})

	t.Run("Test Status LessThan", func(t *testing.T) {
		var args = []string{
			"http://localhost:40000/pretty",
			"--status",
			"--lt",
			"201",
		}

		var result = schmokin.Run(args)
		assert.True(t, result.Success())
	})

	t.Run("Test Status LessThanOrEqual Less", func(t *testing.T) {
		var args = []string{
			"http://localhost:40000/pretty",
			"--status",
			"--lte",
			"201",
		}

		var result = schmokin.Run(args)
		assert.True(t, result.Success())
	})

	t.Run("Test Status LessThanOrEqual Equal", func(t *testing.T) {
		var args = []string{
			"http://localhost:40000/pretty",
			"--status",
			"--lte",
			"200",
		}

		var result = schmokin.Run(args)
		assert.True(t, result.Success())
	})

	t.Run("Test Body Contains", func(t *testing.T) {
		var args = []string{
			"http://localhost:40000/pretty",
			"--res-body",
			"--co",
			"OK",
		}

		var result = schmokin.Run(args)
		assert.True(t, result.Success())
	})

	t.Run("--resp-header", func(t *testing.T) {
		var args = []string{
			"http://localhost:40000/pretty",
			"--res-header",
			"X-FU",
			"--eq",
			"BAR",
		}

		var result = schmokin.Run(args)
		assert.True(t, result.Success())
	})

	t.Run("-- -X POST -d 'UP'", func(t *testing.T) {
		var args = []string{
			"http://localhost:40000/pretty",
			"--res-body",
			"--eq",
			"Method: POST Body: UP",
			"--",
			"-X",
			"POST",
			"-d",
			"UP",
		}

		var result = schmokin.Run(args)
		assert.True(t, result.Success())
	})

	t.Run("--export", func(t *testing.T) {
		var args = []string{
			"http://localhost:40000/pretty",
			"--res-body",
			"--export",
			"TheBody",
			"--",
			"-X",
			"POST",
			"-d",
			"UP",
		}

		schmokin.Run(args)
		args = []string{
			"http://localhost:40000/echo",
			"--res-body",
			"--eq",
			"TheBody",
			"--",
			"-X",
			"POST",
			"-d",
			"$TheBody",
		}
		var result = schmokin.Run(args)
		assert.True(t, result.Success())
	})

	t.Run("-f", func(t *testing.T) {
		tmpfile, err := ioutil.TempFile("", "example")
		if err != nil {
			log.Fatal(err)
		}

		defer os.Remove(tmpfile.Name()) // clean up
		tmpfile.WriteString("http://localhost:40000/pretty --status --eq 200\n")
		tmpfile.WriteString("http://localhost:40000/pretty --status --eq 200\n")

		var args = []string{
			"-f",
			tmpfile.Name(),
		}

		var result = schmokin.Run(args)
		assert.True(t, result.Success())

		if err := tmpfile.Close(); err != nil {
			log.Fatal(err)
		}
	})
}
