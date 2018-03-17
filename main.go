package main

import (
	"flag"
	"fmt"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Context struct {
	ListenHost     string `default:"0.0.0.0" split_words:"true" desc:"Host to listen on"`
	ListenPort     uint   `default:"80" split_words:"true" desc:"Port to listen on"`
	RedirectHost   string `default:"" split_words:"true" desc:"Host to redirect to. Empty hosts mean the host from HTTP request will be used."`
	RedirectPort   uint   `default:"443" split_words:"true" desc:"Port to redirect to"`
	RedirectStatus int    `default:"301" split_words:"true"`
	RedirectScheme string `default:"https" split_words:"true"`
	LogFormat      string `default:"txt" split_words:"true" desc:"Log format. Allowed values are 'txt' and 'json'"`
	LogLevel       string `default:"info" split_words:"true"`
}

func (c Context) redirectHandler(w http.ResponseWriter, req *http.Request) {
	host, err := getHostName(c.RedirectHost, req)

	if err != nil {
		log.Fatal(err, req.Host)
		os.Exit(1)
	}

	u := url.URL(*req.URL)

	u.Host = ""
	u.Scheme = ""

	src := fmt.Sprintf("http://%s%s", req.Host, u.String())
	target := fmt.Sprintf("%s://%s:%d%s", c.RedirectScheme, host, c.RedirectPort, u.String())

	defer log.Debugf("%s %s -> %s %d", req.Method, src, target, c.RedirectStatus)

	http.Redirect(w, req, target, c.RedirectStatus)
}

func getHostName(host string, req *http.Request) (string, error) {
	if "" != host {
		return host, nil
	}

	if strings.IndexByte(req.Host, ':') < 1 {
		return req.Host, nil
	}

	h, _, err := net.SplitHostPort(req.Host)

	if err != nil {
		return "", err
	}

	return h, nil
}

func initLog(c Context) error {
	if "json" == c.LogFormat {
		log.SetFormatter(&log.JSONFormatter{})
	}

	level, err := log.ParseLevel(c.LogLevel)

	if err == nil {
		log.SetLevel(level)
	}

	return err
}

func main() {
	var c Context

	h := flag.Bool("h", false, "Print this help")
	d := flag.Bool("d", false, "Dump config values")
	flag.Parse()

	if *h {
		flag.Usage()
		fmt.Print("\n")
		envconfig.Usage("", &c)
		return
	}

	err := envconfig.Process("", &c)
	if err != nil {
		log.Fatal(err.Error())
		os.Exit(1)
	}

	if *d {
		fmt.Printf("%+v\n", c)
		return
	}

	err = initLog(c)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	handler := http.HandlerFunc(c.redirectHandler)

	address := fmt.Sprintf("%s:%d", c.ListenHost, c.ListenPort)
	log.Info(fmt.Sprintf("Listening on %s", address))

	err = http.ListenAndServe(address, handler)

	if err != nil {
		log.Fatal(err)
	}
}
