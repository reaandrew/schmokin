package main

import (
	"testing"
	"net/http"
	"context"
	"log"

	"github.com/stretchr/testify/assert"
)

func Test_Schmokin(t *testing.T){

	
	m := http.NewServeMux()
	s := http.Server{Addr: ":40000", Handler: m}
	ctx, cancel := context.WithCancel(context.Background())
	defer func(){
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

	t.Run("Test_Status_Equals", func (t *testing.T) {
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

	t.Run("Test_Status_NotEquals", func (t *testing.T) {
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
}
