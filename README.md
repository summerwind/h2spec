# h2spec

h2spec is a conformance testing tool for HTTP/2 implementation.  
This tool is compliant with [RFC 7540 (HTTP/2)](http://www.rfc-editor.org/rfc/rfc7540.txt).

## Install

Go to the [releases page](https://github.com/summerwind/h2spec/releases), find the version you want, and download the zip file.

## Build

1. Make sure you have go 1.5 and set GOPATH appropriately
2. Run `go get github.com/summerwind/h2spec/cmd/h2spec`

It is also possible to build specific version.

1. Make sure you have go 1.5 and set GOPATH appropriately
2. Run `go get gopkg.in/summerwind/h2spec.v1/cmd/h2spec`

## Usage

```
$ h2spec --help
Usage: h2spec [OPTIONS]

Options:
  -p:        Target port. (Default: 80 or 443)
  -h:        Target host. (Default: 127.0.0.1)
  -t:        Connect over TLS. (Default: false)
  -k:        Don't verify server's certificate. (Default: false)
  -o:        Maximum time allowed for test. (Default: 2)
  -s:        Section number on which to run the test. (Example: -s 6.1 -s 6.2)
  -S:        Run the test cases marked as "strict".
  -j:        Creates report also in JUnit format into specified file.
  --version: Display version information and exit.
  --help:    Display this help and exit.
```

## Screenshot

![Sceenshot](https://cloud.githubusercontent.com/assets/230145/6203647/bb15df9e-b56f-11e4-864e-fc63ac0743fb.png)

