package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

const URL = "http://my-prod-service.com"

func main() {
	startServer(execReq, URL)
}

func startServer(exec execRegFn, url string) {
	r := mux.NewRouter()

	r.PathPrefix("/").HandlerFunc(proxy(exec, url))

	srv := &http.Server{
		Handler:      r,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	log.Println("Proxy service successfully started.")
	log.Fatalf("Proxy service error: , %s\n", srv.ListenAndServe())
}

func proxy(exec execRegFn, url string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		u := fmt.Sprintf("%s%s", url, r.URL)
		log.Println(u)
		req, err := http.NewRequest(r.Method, u, r.Body)
		if err != nil {
			proxyErr(w, err)
			return
		}

		req.Header = r.Header
		cl := http.Client{}
		resp, err := exec(cl, req)
		if err != nil {
			proxyErr(w, err)
			return
		}
		defer resp.Body.Close()

		for name, values := range resp.Header {
			w.Header().Set(name, strings.Join(values, ", "))
		}
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	}
}

func proxyErr(w http.ResponseWriter, err error) {
	log.Printf("proxy service error: %v\n", err)
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("proxy service error"))
}

type execRegFn func(client http.Client, req *http.Request) (*http.Response, error)

func execReq(client http.Client, req *http.Request) (*http.Response, error) {
	return client.Do(req)
}
