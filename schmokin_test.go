package main_test

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	schmokin "github.com/reaandrew/schmokin"
)

func Test_Schmokin(t *testing.T) {
	server := CreateTestServer()
	defer server.Stop()
	go server.Start()

	time.Sleep(2 * time.Second)

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
			"http://localhost:40000/echo",
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
			"UP",
			"--",
			"-X",
			"POST",
			"-d",
			"$TheBody",
		}
		var result = schmokin.Run(args)
		assert.True(t, result.Success())
	})

	t.Run("--export with a file", func(t *testing.T) {
		testFilePath := "/tmp/data"
		schmokin.WriteFile(testFilePath, []byte("UP"))
		var args = []string{
			"http://localhost:40000/echo",
			"--res-body",
			"--export",
			"TheBody",
			"--",
			"-X",
			"POST",
			"-d",
			"@" + testFilePath,
		}

		schmokin.Run(args)
		args = []string{
			"http://localhost:40000/echo",
			"--res-body",
			"--eq",
			"UP",
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
		fmt.Println(result)
		assert.True(t, result.Success())

		if err := tmpfile.Close(); err != nil {
			log.Fatal(err)
		}
	})

	t.Run("Invoke without any arguments", func(t *testing.T) {
		var result = schmokin.Run([]string{})
		assert.NotNil(t, result.Error)
	})

	t.Run("Invoke without the --", func(t *testing.T) {
		args := []string{
			"http://localhost:40000/echo_method",
			"-X",
			"POST",
			"--res-body",
			"--eq",
			"POST",
		}
		var result = schmokin.Run(args)
		assert.Nil(t, result.Error)
		assert.True(t, result.Success())
	})
}
