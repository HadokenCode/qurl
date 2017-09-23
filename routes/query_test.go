// Copyright 2017 The qurl Authors. All rights reserved.
package routes

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Fooo", "bar")
			fmt.Fprintf(w, `
				<html>
					<head>
						<title>Page Title</title>
					</head>
					<body>
						<div class="classname">selector class content</div>
						<div id="idname">selector id content</div>
					</body>
				</html>
			`)
		})
		http.ListenAndServe(":6060", nil)
	}()
	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestQuery(t *testing.T) {
	req, err := http.NewRequest("GET", "/q?url=http://localhost:6060", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(Query)

	handler.ServeHTTP(rec, req)

	if status := rec.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"url":"http://localhost:6060","status":200}`
	if rec.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rec.Body.String(), expected)
	}
}

func TestQueryInvalidURL(t *testing.T) {
	req, err := http.NewRequest("GET", "/q?url=invalidurl", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(Query)

	handler.ServeHTTP(rec, req)

	// Check the status code is what we expect.
	if status := rec.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	// Check the response body is what we expect.
	expected := "INVALID_URL\n"
	if rec.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rec.Body.String(), expected)
	}
}

func TestFailFetchURL(t *testing.T) {
	req, err := http.NewRequest("GET", "/q?url=http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(Query)

	handler.ServeHTTP(rec, req)

	// Check the status code is what we expect.
	if status := rec.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}

	// Check the response body is what we expect.
	expected := "INTERNAL_ERROR\n"
	if rec.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rec.Body.String(), expected)
	}
}
