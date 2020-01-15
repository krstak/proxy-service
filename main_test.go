package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/krstak/testify"
)

func TestProxy(t *testing.T) {
	var remoteURL string
	dummyResp := func(client http.Client, req *http.Request) (*http.Response, error) {
		remoteURL = fmt.Sprintf("%s", req.URL)
		return &http.Response{Header: dummyHeader(), StatusCode: 201, Body: ioutil.NopCloser(bytes.NewReader([]byte("test response")))}, nil
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/api/articles", nil)
	proxy(dummyResp, "http://prod.com")(w, r)

	testify.Equal(t)("http://prod.com/api/articles", remoteURL)
	testify.Equal(t)("test response", w.Body.String())
	testify.Equal(t)(dummyHeader(), w.HeaderMap)
	testify.Equal(t)(201, w.Code)
}

func TestProxy_DummyServerReturnsError(t *testing.T) {
	dummyResp := func(client http.Client, req *http.Request) (*http.Response, error) {
		return &http.Response{}, errors.New("error from dummy server")
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/api/articles", nil)
	proxy(dummyResp, "http://prod.com")(w, r)

	testify.Equal(t)("proxy service error", w.Body.String())
	testify.Equal(t)(make(http.Header), w.HeaderMap)
	testify.Equal(t)(500, w.Code)
}

func dummyHeader() http.Header {
	header := make(http.Header)
	header.Add("key", "val")
	return header
}
