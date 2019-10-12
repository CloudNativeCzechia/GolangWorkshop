package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRandString(t *testing.T) {
	// rand string should return different string after each call
	var oldString string
	for i := 0; i < 3; i++ {
		currentString := randString(5)
		if currentString == oldString {
			t.Errorf("randString should always return new string, %s vs %s", oldString, currentString)
		}
		oldString = currentString
	}
}

func TestHandleShorten(t *testing.T) {
	url := "http://seznam.cz"

	req, err := http.NewRequest(
		http.MethodGet,
		"http://localhost:8080/shorten?url="+url,
		nil,
	)
	if err != nil {
		t.Fatalf("Error creating request, %s", err.Error())
	}

	rec := httptest.NewRecorder()
	handleShorten(rec, req)

	if rec.Code != 200 {
		t.Errorf("Processing returned and error: %d", rec.Code)
	}

	if strings.Contains(rec.Body.String(), url) {
		t.Error("Shortened url is the same")
	}
}

func TestHandleRedirect(t *testing.T) {
	url := "http://seznam.cz"
	shortURL := randString(5)
	// mock the work of shortenHandler
	urlsMap[shortURL] = url

	req, err := http.NewRequest(
		http.MethodGet,
		"http://localhost:8080/r?q="+shortURL,
		nil,
	)
	if err != nil {
		t.Fatalf("Error creating request, %s", err.Error())
	}

	rec := httptest.NewRecorder()
	handleRedirect(rec, req)
	if rec.Code != 303 {
		t.Errorf("Processing returned and error: %d", rec.Code)
	}

	if rec.HeaderMap.Get("location") != url {
		t.Errorf("Redirecting to the wrong location, %s", rec.HeaderMap.Get("location"))
	}
}

func TestMain(t *testing.T) {

}
