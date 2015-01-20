# h2spec

h2spec is a conformance test tool for HTTP/2 server.  
This tool supports [draft-ietf-httpbis-http2-16](http://tools.ietf.org/html/draft-ietf-httpbis-http2-16).

## Install

Go to the [releases page](https://github.com/summerwind/h2spec/releases), find the version you want, and download the zip file.

## Build

1. Make sure you have go 1.4 and set GOPATH appropriately
2. Clone this repo and move to the cloned repo directory
3. Run `go get github.com/bradfitz/http2`
4. Run `go build cmd/h2spec.go`

## Usage

```
$ h2spec --help
Usage: h2spec [OPTIONS]

Options:
  -p:     Target port. (Default: 80)
  -h:     Target host. (Default: 127.0.0.1)
  -t:     Connect over TLS. (Default: false)
  -k:     Don't verify server's certificate. (Default: false)
  -s:     Section number on which to run the test. (Example: -s 6.1 -s 6.2)
  --help: Display this help and exit.
```

## Screenshot

![Sceenshot](https://cloud.githubusercontent.com/assets/230145/5282230/c267a6b8-7b47-11e4-8949-2121d8921382.png)

