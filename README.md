# http-redirector

[![Build Status](https://api.travis-ci.org/pinepain/http-redirector.svg?branch=master)](https://travis-ci.org/pinepain/http-redirector)
[![codecov](https://codecov.io/gh/pinepain/http-redirector/branch/master/graph/badge.svg)](https://codecov.io/gh/pinepain/http-redirector)
[![Go Report Card](https://goreportcard.com/badge/github.com/pinepain/http-redirector)](https://goreportcard.com/report/github.com/pinepain/http-redirector)

Simple `http` to anything redirector. By default it is `http` -> `https` redirector.

## Rationale

Often we may have HTTPS-only, self-sufficient services, however we may have to redirect traffic from HTTP for
convenience. The most common case is to have service behind AWS ELB with TLS termination. While you still set
up separate HTTP web server like nginx, apache or whatever, there is always a Unix way to this - to have a program
that does one thing well. This is what `http-redictor` for.

## Key features

 - it works;
 - configurable (it could even write logs in json);
 - written in golang and even has few tests;
 - has own Docker image;
 - uses hyped technologies;
 - it's my first service in golang.

## Usage

By default `http-redirector` listens on port `80` and redirects all `http` traffic to the `https`
with HTTP 301 Moved Permanently code.

The recommended way to use is to use `pinepain/http-redirector` docker image, e.g.:
`$ docker run -p 8080:80 -e LOG_LEVEL=DEBUG pinepain/http-redirector`
 
then in different console let's send http query to port `8080`: 

```
$ curl -v 'http://localhost:8080/foo/bar?test=me'
*   Trying ::1...
* TCP_NODELAY set
* Connected to localhost (::1) port 8080 (#0)
> GET /foo/bar?test=me HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.54.0
> Accept: */*
> 
< HTTP/1.1 301 Moved Permanently
< Content-Type: text/html; charset=utf-8
< Location: https://localhost/foo/bar?test=me
< Date: Sat, 17 Mar 2018 19:18:06 GMT
< Content-Length: 72
< 
<a href="https://localhost/foo/bar?test=me">Moved Permanently</a>.

* Connection #0 to host localhost left intact
```

## Configuration

`http-redirector` configured with environment variables. To list them, as well as default,
please, run `./http-redirector` with `-h` flag:

```
Usage of ./http-redirector:
  -d	Dump config values
  -h	Print this help

This application is configured via the environment. The following environment
variables can be used:

KEY           TYPE       DEFAULT                 REQUIRED    DESCRIPTION
LISTEN        String     0.0.0.0:80                          Host:port to listen on
REDIRECT      String     https://REQUEST_HOST                Destination to redirect to. You can specify schema, host and port.
                                                             REQUEST_HOST or empty hostname means that request hostname will be used.
STATUS        Integer    301                                 Redirect status code
LOG_FORMAT    String     txt                                 Log format. Allowed values are 'txt' and 'json'
LOG_LEVEL     String     info                                Logs verbosity
```

To see what actual configuration values are, run `./http-redirector` with `-d` flag:

```
export LOG_LEVEL=debug
export LISTEN=8080
export REDIRECT=9090
$ ./http-redirector -d
&{Listen:8080 Redirect:9090 Status:301 LogFormat:txt LogLevel:debug}
```

## License

[http-redirector](https://github.com/pinepain/http-redirector) is licensed under the [MIT license](http://opensource.org/licenses/MIT).
