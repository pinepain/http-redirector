package main

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func responseFor(c *Config, url string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	r, err := http.NewRequest("TEST", url, nil)

	if err != nil {
		logrus.SetOutput(os.Stderr)
		logrus.Fatal(err)
	}

	NewHTTPRedirector(c).ServeHTTP(w, r)

	return w
}

func shouldRedirect(t *testing.T, redirect string, status int, url string, newURL string) {
	c := &Config{}
	c.Redirect = redirect
	c.Status = status

	w := responseFor(c, url)

	if w.Code != status {
		t.Errorf("%v yields invalid response status", url)
	}
	if w.HeaderMap["Location"][0] != newURL {
		t.Errorf("%v should redirect to %v, redirected to %v", url, newURL, w.HeaderMap["Location"][0])
	}
}

func TestRedirect(t *testing.T) {
	// logrus.SetOutput(ioutil.Discard)

	shouldRedirect(t, "https://REQUEST_HOST:443", http.StatusMovedPermanently, "http://foo.com:8080/asd/fdfd?asd=123", "https://foo.com/asd/fdfd?asd=123")
	shouldRedirect(t, ":443", http.StatusMovedPermanently, "http://foo.com/asd/fdfd?asd=123", "https://foo.com/asd/fdfd?asd=123")
	shouldRedirect(t, ":80", http.StatusMovedPermanently, "https://foo.com/asd/fdfd?asd=123", "http://foo.com/asd/fdfd?asd=123")
	shouldRedirect(t, "http://:80", http.StatusMovedPermanently, "https://foo.com/asd/fdfd?asd=123", "http://foo.com/asd/fdfd?asd=123")
	shouldRedirect(t, "https://:443", http.StatusMovedPermanently, "http://foo.com/asd/fdfd?asd=123", "https://foo.com/asd/fdfd?asd=123")
	shouldRedirect(t, "example.com:8080", http.StatusMovedPermanently, "http://foo.com/asd/fdfd?asd=123", "http://example.com:8080/asd/fdfd?asd=123")
	shouldRedirect(t, "example.com:8080/test", http.StatusMovedPermanently, "http://foo.com/asd/fdfd?asd=123", "http://example.com:8080/asd/fdfd?asd=123")
	shouldRedirect(t, "example.com/test", http.StatusMovedPermanently, "http://foo.com/asd/fdfd?asd=123", "http://example.com/asd/fdfd?asd=123")
	shouldRedirect(t, "http://:81", http.StatusFound, "http://127.0.0.1?asd=123", "http://127.0.0.1:81?asd=123")
	shouldRedirect(t, "192.168.0.1:8000", http.StatusFound, "http://127.0.0.1", "http://192.168.0.1:8000")
}
