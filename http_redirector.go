package main

import (
	log "github.com/sirupsen/logrus"
	"net/http"
)

// HTTPRedirector struct
type HTTPRedirector struct {
	Redirect *Redirect
	Status   int
}

func (h *HTTPRedirector) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	u := r.URL

	u.Host = r.Host

	src := "http://" + r.Host + u.String()

	target := h.Redirect.GetHost(r.URL)
	u.Scheme = ""
	u.Host = ""
	target += u.String()

	defer log.WithFields(log.Fields{
		"method": r.Method,
		"src":    src,
		"target": target,
		"status": h.Status,
	}).Debugf("%s %s -> %s %d", r.Method, src, target, h.Status)

	http.Redirect(w, r, target, h.Status)
}

// NewHTTPRedirector - create new HTTPRedirector from config
func NewHTTPRedirector(c *Config) *HTTPRedirector {

	redirector := &HTTPRedirector{
		Redirect: NewRedirect(c.Redirect),
		Status:   c.Status,
	}

	return redirector
}
