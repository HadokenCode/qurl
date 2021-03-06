// Copyright 2017 The qurl Authors. All rights reserved.

package qurl

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

// IRequest ...
type IRequest interface {
	Fetch(url string) (*http.Response, error)
}

// Request represents the call being made to retrieve the contents of an URL.
type Request struct {
}

// Fetch performs an HTTP GET call to anURL and fetch the contents.
func (r *Request) Fetch(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	return client.Do(req)
}

// FakeRequest ...
type FakeRequest struct {
	ExpectedStatusCode      int
	ExpectedBody            string
	ExpectedResponseHeaders http.Header
}

// Fetch performs an HTTP GET call to anURL and fetch the contents.
func (r *FakeRequest) Fetch(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	body := r.ExpectedBody
	resp := &http.Response{
		Status:        http.StatusText(r.ExpectedStatusCode),
		StatusCode:    r.ExpectedStatusCode,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Body:          ioutil.NopCloser(bytes.NewBufferString(body)),
		ContentLength: int64(len(body)),
		Request:       req,
		Header:        r.ExpectedResponseHeaders,
	}
	return resp, nil
}
