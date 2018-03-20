package main

import (
	log "github.com/sirupsen/logrus"
	"net/url"
	"strings"
)

// Redirect struct
type Redirect struct {
	Schema string
	Host   string
	Port   string
}

// NewRedirect - creates new Redirect from redirect URL
func NewRedirect(raw string) *Redirect {
	u := parseRedirectURL(raw)

	u.Scheme = getScheme(u)
	port := getPort(u)
	host := getHost(u)

	return &Redirect{
		Schema: u.Scheme,
		Host:   host,
		Port:   port,
	}
}

// GetHost - get full hostname to redirect URL to
func (r *Redirect) GetHost(u *url.URL) string {
	host := r.Host

	if host == "" {
		host = u.Hostname()
	}

	s := r.Schema + "://" + host

	if r.Port != "" {
		s += ":" + r.Port
	}

	return s
}

func parseRedirectURL(raw string) *url.URL {

	if strings.Index(raw, "://") < 0 {
		redirectURL := &url.URL{}
		redirectURL.Host = raw

		if strings.Index(raw, "/") < 0 {
			return redirectURL
		}

		redirectURL.Host = raw[0:strings.Index(raw, "/")]
		return redirectURL
	}

	redirectURL, err := url.Parse(raw)

	if err != nil {
		log.Fatal(err)
	}

	return redirectURL
}

func getScheme(u *url.URL) string {
	port := u.Port()
	scheme := u.Scheme

	if scheme != "" {
		return scheme
	}
	if "443" == port {
		return "https"
	}

	if port == "80" {
		return "http"
	}

	return "http"
}

func getPort(u *url.URL) string {
	port := u.Port()
	scheme := u.Scheme

	if scheme == "https" && port == "443" {
		return ""
	}

	if scheme == "http" && port == "80" {
		port = ""
	}

	return port
}

func getHost(u *url.URL) string {
	host := u.Hostname()

	if host == "REQUEST_HOST" {
		host = ""
	}

	return host
}
