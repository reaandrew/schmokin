package main_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

type TestServer struct {
	server  *http.Server
	context context.Context
	cancel  context.CancelFunc
}

func CreateTestServer() TestServer {
	m := http.NewServeMux()
	s := http.Server{Addr: ":40000", Handler: m}
	ctx, cancel := context.WithCancel(context.Background())

	m.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			body = []byte("not set")
		}
		w.Write(body)
	})
	m.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		value, _ := strconv.Atoi(r.URL.Query().Get("value"))
		w.WriteHeader(value)
	})
	m.HandleFunc("/json", func(w http.ResponseWriter, r *http.Request) {
		json := `{"name":{"first":"Janet","last":"Prichard"},"age":47}`
		w.Write([]byte(json))
	})
	m.HandleFunc("/echo_method", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(r.Method))
	})
	m.HandleFunc("/redirect", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "http://localhost:40000/echo_method", 301)
	})
	m.HandleFunc("/echo_headers", func(w http.ResponseWriter, r *http.Request) {
		for key, value := range r.Header {
			w.Header().Set(key, strings.Join(value, ","))
		}
		w.Write([]byte(r.Method))
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

	return TestServer{
		server:  &s,
		context: ctx,
		cancel:  cancel,
	}
}

func (instance *TestServer) Start() {
	if err := instance.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func (instance *TestServer) Stop() {
	instance.cancel()
	instance.server.Shutdown(instance.context)
}
