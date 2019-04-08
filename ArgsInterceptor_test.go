package main_test

import (
	"io/ioutil"
	"os"
	"testing"

	schmokin "github.com/reaandrew/schmokin"
	"github.com/stretchr/testify/assert"
)

func Test_ArgInterceptor(t *testing.T) {

	t.Run("-d with data", func(t *testing.T) {
		args := []string{"-d", "$value"}

		state := schmokin.State{
			"value": "boo",
		}

		interceptor := schmokin.CreateArgsInterceptor(state)

		newArgs := interceptor.Intercept(args)

		assert.Equal(t, []string{
			"-d",
			"boo",
		}, newArgs)
	})

	t.Run("-d with file", func(t *testing.T) {
		data := []byte("hello\n$value\n")
		f, _ := os.Create("./data")
		defer func() {
			f.Close()
			os.Remove(f.Name())
		}()
		f.Write(data)

		args := []string{"-d", "@data"}

		state := schmokin.State{
			"value": "boo",
		}

		interceptor := schmokin.CreateArgsInterceptor(state)

		newArgs := interceptor.Intercept(args)

		dat, _ := ioutil.ReadFile("schmokin.payload")

		assert.Equal(t, string(dat), "hello\nboo\n")
		assert.Equal(t, []string{
			"-d",
			"@schmokin.payload",
		}, newArgs)
	})
}
