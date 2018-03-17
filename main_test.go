package main

import (
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func responseFor(ctxt Context, url string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("TEST", url, nil)
	ctxt.redirectHandler(w, r)
	return w
}

func shouldRedirect(t *testing.T, host string, port uint, status int, scheme string, url string, newUrl string) {
	ctxt := Context{}
	ctxt.RedirectHost = host
	ctxt.RedirectPort = port
	ctxt.RedirectStatus = status
	ctxt.RedirectScheme = scheme

	w := responseFor(ctxt, url)
	if w.Code != status {
		t.Errorf("%v should yirld a 301 response", url)
	}
	if w.HeaderMap["Location"][0] != newUrl {
		t.Errorf("%v should redirect to %v, redirected to %v", url, newUrl, w.HeaderMap["Location"][0])
	}
}

func TestRedirect(t *testing.T) {
	logrus.SetOutput(ioutil.Discard)

	shouldRedirect(t, "", 443, http.StatusMovedPermanently, "https", "http://foo.com:8080/asd/fdfd?asd=123", "https://foo.com:443/asd/fdfd?asd=123")
	shouldRedirect(t, "", 443, http.StatusMovedPermanently, "https", "http://foo.com/asd/fdfd?asd=123", "https://foo.com:443/asd/fdfd?asd=123")
	shouldRedirect(t, "example.com", 8080, http.StatusMovedPermanently, "http", "http://foo.com/asd/fdfd?asd=123", "http://example.com:8080/asd/fdfd?asd=123")
	shouldRedirect(t, "", 81, http.StatusFound, "http", "http://127.0.0.1?asd=123", "http://127.0.0.1:81?asd=123")
	shouldRedirect(t, "192.168.0.1", 8000, http.StatusFound, "http", "http://127.0.0.1", "http://192.168.0.1:8000")
}
