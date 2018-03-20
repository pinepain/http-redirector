package main

import (
	"flag"
	"fmt"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// Config struct
type Config struct {
	Listen    string `default:"0.0.0.0:80" split_words:"true" desc:"Host:port to listen on"`
	Redirect  string `default:"https://REQUEST_HOST" split_words:"true" desc:"Destination to redirect to. You can specify schema, host and port. REQUEST_HOST or empty hostname means that request hostname will be used."`
	Status    int    `default:"301" split_words:"true" desc:"Redirect status code"`
	LogFormat string `default:"txt" split_words:"true" desc:"Log format. Allowed values are 'txt' and 'json'"`
	LogLevel  string `default:"info" split_words:"true" desc:"Logs verbosity"`
}

func initLog(c *Config) {
	if "json" == c.LogFormat {
		log.SetFormatter(&log.JSONFormatter{})
	}

	level, err := log.ParseLevel(c.LogLevel)

	if err != nil {

		log.Fatal(err)
	}

	log.SetLevel(level)
}

func main() {
	c := &Config{}

	h := flag.Bool("h", false, "Print this help")
	d := flag.Bool("d", false, "Dump config values")
	flag.Parse()

	if *h {
		flag.Usage()
		fmt.Print("\n")
		envconfig.Usage("", c)
		return
	}

	err := envconfig.Process("", c)
	if err != nil {
		log.Fatal(err.Error())
	}

	if *d {
		fmt.Printf("%+v\n", c)
		return
	}

	initLog(c)

	r := NewHTTPRedirector(c)

	log.Info(fmt.Sprintf("Listening on %s", c.Listen))
	log.Fatal(http.ListenAndServe(c.Listen, r))
}
